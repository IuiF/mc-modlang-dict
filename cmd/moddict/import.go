package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/internal/jar"
	"github.com/iuif/minecraft-mod-dictionary/internal/parser"
	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ExitOnError)

	var (
		dbPath   = fs.String("db", "moddict.db", "Database file path")
		workDir  = fs.String("work", "workspace/temp", "Working directory for extraction")
		jarPath  = fs.String("jar", "", "Path to mod JAR file (required)")
		langCode = fs.String("lang", "en_us", "Source language code")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict import [options]

Import translations from a Minecraft mod JAR file.
Reuses existing sources if mod_id + key + source_text matches.
Sets the imported version as the default version.

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  moddict import -jar create-1.20.1-0.5.1.jar
  moddict import -jar mods/botania.jar -db translations.db
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *jarPath == "" {
		fs.Usage()
		return fmt.Errorf("JAR file path is required")
	}

	// Verify JAR exists
	if _, err := os.Stat(*jarPath); os.IsNotExist(err) {
		return fmt.Errorf("JAR file not found: %s", *jarPath)
	}

	// Create work directory
	if err := os.MkdirAll(*workDir, 0755); err != nil {
		return fmt.Errorf("failed to create work directory: %w", err)
	}

	// Open database
	repo, err := database.NewRepository(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	// Run migrations
	if err := repo.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Extract JAR
	fmt.Printf("Extracting %s...\n", filepath.Base(*jarPath))
	extractor := jar.NewExtractor()
	extractDir := filepath.Join(*workDir, filepath.Base(*jarPath))

	result, err := extractor.Extract(*jarPath, extractDir)
	if err != nil {
		return fmt.Errorf("failed to extract JAR: %w", err)
	}

	fmt.Printf("Detected mod: %s (%s)\n", result.DisplayName, result.ModID)
	fmt.Printf("Loader: %s, Version: %s\n", result.Loader, result.Version)
	fmt.Printf("Found %d lang files\n", len(result.LangFiles))

	ctx := context.Background()

	// Save mod info
	mod := &models.Mod{
		ID:          result.ModID,
		DisplayName: result.DisplayName,
		Author:      joinAuthors(result.Authors),
		Description: result.Description,
	}

	if err := repo.SaveMod(ctx, mod); err != nil {
		return fmt.Errorf("failed to save mod: %w", err)
	}

	// Get or create version (reuses existing if mod_id + version matches)
	modVersion, versionCreated, err := repo.GetOrCreateVersion(ctx, result.ModID, result.Version, result.Loader)
	if err != nil {
		return fmt.Errorf("failed to get/create version: %w", err)
	}

	if versionCreated {
		fmt.Printf("Created new version: %s\n", result.Version)
	} else {
		fmt.Printf("Reusing existing version: %s (ID=%d)\n", result.Version, modVersion.ID)
	}

	// Set this version as default
	if err := repo.SetDefaultVersion(ctx, modVersion.ID); err != nil {
		return fmt.Errorf("failed to set default version: %w", err)
	}
	modVersion.IsDefault = true // Sync local state with DB

	// Parse and import lang files
	jsonParser := parser.NewJSONLangParser()
	legacyParser := parser.NewLegacyLangParser()
	var totalKeys, reusedSources, newSources, reusedTranslations, newTranslations, officialTranslations int

	// First pass: Import source language (en_us) and create sources
	sourceEntries := make(map[string]interfaces.ParsedEntry) // key -> entry
	for _, langFile := range result.LangFiles {
		if !isSourceLangFile(langFile, *langCode) {
			continue
		}

		content, err := os.ReadFile(langFile)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", langFile, err)
			continue
		}

		var langParser interfaces.Parser
		if isLegacyLangFile(langFile) {
			langParser = legacyParser
		} else {
			langParser = jsonParser
		}

		entries, err := langParser.Parse(content)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", langFile, err)
			continue
		}

		for _, entry := range entries {
			sourceEntries[entry.Key] = entry
		}

		totalKeys += len(entries)
		fmt.Printf("Processed %d keys from %s\n", len(entries), filepath.Base(langFile))
	}

	// Second pass: Load official Japanese translations if available
	jaTranslations := make(map[string]string) // key -> japanese text
	for _, langFile := range result.LangFiles {
		if !isTargetLangFile(langFile, "ja_jp") {
			continue
		}

		content, err := os.ReadFile(langFile)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", langFile, err)
			continue
		}

		var langParser interfaces.Parser
		if isLegacyLangFile(langFile) {
			langParser = legacyParser
		} else {
			langParser = jsonParser
		}

		entries, err := langParser.Parse(content)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", langFile, err)
			continue
		}

		for _, entry := range entries {
			jaTranslations[entry.Key] = entry.Text
		}

		fmt.Printf("Found %d official Japanese translations from %s\n", len(entries), filepath.Base(langFile))
	}

	// Third pass: Create sources and translations
	for key, entry := range sourceEntries {
		// Get or create source (reuses if mod_id + key + source_text matches)
		source, created, err := repo.GetOrCreateSource(ctx, result.ModID, entry.Key, entry.Text, *langCode)
		if err != nil {
			return fmt.Errorf("failed to get/create source: %w", err)
		}

		if created {
			newSources++
		} else {
			reusedSources++
		}

		// Link source to this version
		if err := repo.LinkSourceToVersion(ctx, source.ID, modVersion.ID); err != nil {
			return fmt.Errorf("failed to link source to version: %w", err)
		}

		// Check if translation already exists for this source
		existingTrans, err := repo.GetTranslationForSource(ctx, source.ID)
		if err != nil {
			return fmt.Errorf("failed to check existing translation: %w", err)
		}

		if existingTrans != nil {
			// Translation exists - check if we should update with official translation
			if jaText, hasOfficial := jaTranslations[key]; hasOfficial {
				// Update with official translation if current is pending or empty
				if existingTrans.Status == models.StatusPending || existingTrans.TargetText == nil || *existingTrans.TargetText == "" {
					existingTrans.TargetText = &jaText
					existingTrans.Status = models.StatusOfficial
					if err := repo.SaveTranslation(ctx, existingTrans); err != nil {
						return fmt.Errorf("failed to update translation: %w", err)
					}
					officialTranslations++
				}
			}
			reusedTranslations++
		} else {
			// Create new translation
			trans := &models.Translation{
				SourceID:   source.ID,
				TargetLang: "ja_jp",
				Tags:       entry.Tags,
			}

			// Check if official Japanese translation exists
			if jaText, hasOfficial := jaTranslations[key]; hasOfficial {
				trans.TargetText = &jaText
				trans.Status = models.StatusOfficial
				officialTranslations++
			} else {
				trans.Status = models.StatusPending
			}

			if err := repo.SaveTranslation(ctx, trans); err != nil {
				return fmt.Errorf("failed to save translation: %w", err)
			}
			newTranslations++
		}
	}

	// Update version stats
	modVersion.Stats.TotalKeys = totalKeys
	if err := repo.SaveVersion(ctx, modVersion); err != nil {
		return fmt.Errorf("failed to update version stats: %w", err)
	}

	fmt.Printf("\nImport complete:\n")
	fmt.Printf("  Total keys: %d\n", totalKeys)
	fmt.Printf("  Sources: %d reused, %d new\n", reusedSources, newSources)
	fmt.Printf("  Translations: %d reused, %d new\n", reusedTranslations, newTranslations)
	if officialTranslations > 0 {
		fmt.Printf("  Official Japanese: %d keys\n", officialTranslations)
	}
	return nil
}

func joinAuthors(authors []string) string {
	if len(authors) == 0 {
		return ""
	}
	if len(authors) == 1 {
		return authors[0]
	}

	result := authors[0]
	for i := 1; i < len(authors); i++ {
		result += ", " + authors[i]
	}
	return result
}

func isSourceLangFile(path, langCode string) bool {
	base := filepath.Base(path)
	// Support both .json (1.13+) and .lang (1.12.2 and earlier) formats
	// Use case-insensitive comparison for language code
	baseLower := strings.ToLower(base)
	langCodeLower := strings.ToLower(langCode)
	return baseLower == langCodeLower+".json" || baseLower == langCodeLower+".lang"
}

func isTargetLangFile(path, langCode string) bool {
	base := filepath.Base(path)
	// Support both .json (1.13+) and .lang (1.12.2 and earlier) formats
	// Use case-insensitive comparison for language code
	baseLower := strings.ToLower(base)
	langCodeLower := strings.ToLower(langCode)
	return baseLower == langCodeLower+".json" || baseLower == langCodeLower+".lang"
}

func isLegacyLangFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".lang")
}
