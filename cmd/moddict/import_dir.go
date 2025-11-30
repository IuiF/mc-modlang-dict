package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runImportDir(args []string) error {
	fs := flag.NewFlagSet("import-dir", flag.ExitOnError)

	var (
		dbPath    = fs.String("db", "moddict.db", "Database file path")
		modID     = fs.String("mod", "", "Mod ID (required)")
		version   = fs.String("version", "", "Mod version (required)")
		mcVersion = fs.String("mc", "", "Minecraft version (required)")
		langDir   = fs.String("lang", "", "Path to lang directory containing en_us.json")
		patchDir  = fs.String("patchouli", "", "Path to patchouli entries directory")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict import-dir [options]

Import translations from a directory (cloned repository).

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  moddict import-dir -mod bloodmagic -version 3.0.0 -mc 1.16.3 \
    -lang ./repo/src/main/resources/assets/bloodmagic/lang \
    -patchouli ./repo/src/main/resources/data/bloodmagic/patchouli_books/guide/en_us/entries
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *modID == "" || *version == "" || *mcVersion == "" {
		fs.Usage()
		return fmt.Errorf("mod, version, and mc are required")
	}

	if *langDir == "" && *patchDir == "" {
		return fmt.Errorf("at least one of -lang or -patchouli is required")
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

	ctx := context.Background()

	// Save or get mod
	mod := &models.Mod{
		ID:          *modID,
		DisplayName: *modID,
	}
	if err := repo.SaveMod(ctx, mod); err != nil {
		return fmt.Errorf("failed to save mod: %w", err)
	}

	// Create version
	modVersion := &models.ModVersion{
		ModID:     *modID,
		Version:   *version,
		MCVersion: *mcVersion,
	}
	if err := repo.SaveVersion(ctx, modVersion); err != nil {
		return fmt.Errorf("failed to save version: %w", err)
	}

	var totalKeys int

	// Import lang file
	if *langDir != "" {
		langFile := filepath.Join(*langDir, "en_us.json")
		count, err := importLangFile(ctx, repo, *modID, modVersion.ID, langFile)
		if err != nil {
			return fmt.Errorf("failed to import lang file: %w", err)
		}
		totalKeys += count
		fmt.Printf("Imported %d keys from lang file\n", count)
	}

	// Import patchouli entries
	if *patchDir != "" {
		count, err := importPatchouliEntries(ctx, repo, *modID, modVersion.ID, *patchDir)
		if err != nil {
			return fmt.Errorf("failed to import patchouli entries: %w", err)
		}
		totalKeys += count
		fmt.Printf("Imported %d keys from patchouli entries\n", count)
	}

	// Update version stats
	modVersion.Stats.TotalKeys = totalKeys
	if err := repo.SaveVersion(ctx, modVersion); err != nil {
		return fmt.Errorf("failed to update version stats: %w", err)
	}

	fmt.Printf("\nImport complete: %d total keys\n", totalKeys)
	return nil
}

func importLangFile(ctx context.Context, repo *database.Repository, modID string, versionID int64, langFile string) (int, error) {
	content, err := os.ReadFile(langFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", langFile, err)
	}

	var langData map[string]string
	if err := json.Unmarshal(content, &langData); err != nil {
		return 0, fmt.Errorf("failed to parse %s: %w", langFile, err)
	}

	count := 0
	for key, text := range langData {
		if err := importSingleEntry(ctx, repo, modID, versionID, key, text); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func importPatchouliEntries(ctx context.Context, repo *database.Repository, modID string, versionID int64, entriesDir string) (int, error) {
	var totalKeys int

	err := filepath.Walk(entriesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Warning: failed to read %s: %v\n", path, err)
			return nil
		}

		var entry map[string]interface{}
		if err := json.Unmarshal(content, &entry); err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", path, err)
			return nil
		}

		// Get relative path for key prefix
		relPath, _ := filepath.Rel(entriesDir, path)
		keyPrefix := "patchouli:" + strings.TrimSuffix(relPath, ".json")

		count, err := importPatchouliEntry(ctx, repo, modID, versionID, keyPrefix, entry)
		if err != nil {
			return fmt.Errorf("failed to import %s: %w", path, err)
		}
		totalKeys += count

		return nil
	})

	return totalKeys, err
}

func importPatchouliEntry(ctx context.Context, repo *database.Repository, modID string, versionID int64, keyPrefix string, entry map[string]interface{}) (int, error) {
	count := 0

	// Entry name
	if name, ok := entry["name"].(string); ok {
		if err := importSingleEntry(ctx, repo, modID, versionID, keyPrefix+".name", name); err != nil {
			return count, err
		}
		count++
	}

	// Pages
	if pages, ok := entry["pages"].([]interface{}); ok {
		for i, page := range pages {
			if pageMap, ok := page.(map[string]interface{}); ok {
				// Title
				if title, ok := pageMap["title"].(string); ok {
					if err := importSingleEntry(ctx, repo, modID, versionID, fmt.Sprintf("%s.page%d.title", keyPrefix, i), title); err != nil {
						return count, err
					}
					count++
				}
				// Text
				if text, ok := pageMap["text"].(string); ok {
					if err := importSingleEntry(ctx, repo, modID, versionID, fmt.Sprintf("%s.page%d.text", keyPrefix, i), text); err != nil {
						return count, err
					}
					count++
				}
			}
		}
	}

	return count, nil
}

func importSingleEntry(ctx context.Context, repo *database.Repository, modID string, versionID int64, key, text string) error {
	// Create or get source
	source := &models.TranslationSource{
		ModID:      modID,
		Key:        key,
		SourceText: text,
		SourceLang: "en_us",
		IsCurrent:  true,
	}
	if err := repo.SaveSource(ctx, source); err != nil {
		return fmt.Errorf("failed to save source: %w", err)
	}

	// Link source to version
	if err := repo.LinkSourceToVersion(ctx, source.ID, versionID); err != nil {
		return fmt.Errorf("failed to link source to version: %w", err)
	}

	// Create translation (pending) - include both new and legacy fields
	trans := &models.Translation{
		SourceID:     source.ID,
		ModVersionID: versionID,
		Key:          key,
		SourceText:   text,
		SourceLang:   "en_us",
		TargetLang:   "ja_jp",
		Status:       models.StatusPending,
	}
	if err := repo.SaveTranslation(ctx, trans); err != nil {
		return fmt.Errorf("failed to save translation: %w", err)
	}

	return nil
}
