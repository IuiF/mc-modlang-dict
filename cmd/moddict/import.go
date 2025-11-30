package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/internal/jar"
	"github.com/iuif/minecraft-mod-dictionary/internal/parser"
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

	// Create version (will be set as default)
	modVersion := &models.ModVersion{
		ModID:     result.ModID,
		Version:   result.Version,
		Loader:    result.Loader,
		IsDefault: false, // Will be set after creation
	}

	if err := repo.SaveVersion(ctx, modVersion); err != nil {
		return fmt.Errorf("failed to save version: %w", err)
	}

	// Set this version as default
	if err := repo.SetDefaultVersion(ctx, modVersion.ID); err != nil {
		return fmt.Errorf("failed to set default version: %w", err)
	}

	// Parse and import lang files
	jsonParser := parser.NewJSONLangParser()
	var totalKeys, reusedSources, newSources, reusedTranslations, newTranslations int

	for _, langFile := range result.LangFiles {
		// Only process source language file
		if !isSourceLangFile(langFile, *langCode) {
			continue
		}

		content, err := os.ReadFile(langFile)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", langFile, err)
			continue
		}

		entries, err := jsonParser.Parse(content)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", langFile, err)
			continue
		}

		// Create sources and translations
		for _, entry := range entries {
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
				// Translation exists - reuse it (no need to create new)
				reusedTranslations++
			} else {
				// Create new translation (pending)
				trans := &models.Translation{
					SourceID:   source.ID,
					TargetLang: "ja_jp",
					Status:     models.StatusPending,
					Tags:       entry.Tags,
				}
				if err := repo.SaveTranslation(ctx, trans); err != nil {
					return fmt.Errorf("failed to save translation: %w", err)
				}
				newTranslations++
			}
		}

		totalKeys += len(entries)
		fmt.Printf("Processed %d keys from %s\n", len(entries), filepath.Base(langFile))
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
	return base == langCode+".json"
}
