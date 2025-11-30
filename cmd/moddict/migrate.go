package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runMigrate(args []string) error {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)

	var (
		dbPath = fs.String("db", "moddict.db", "Database file path")
		dryRun = fs.Bool("dry-run", false, "Show what would be migrated without making changes")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict migrate [options]

Migrate existing translations to the new source-based schema.

This command:
1. Creates translation_sources from existing translations
2. Creates source_versions links
3. Updates translations to reference source_id
4. Fixes version information (version, mc_version)

Options:
`)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	repo, err := database.NewRepository(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	// Run migrations first to create new tables
	fmt.Println("Running schema migrations...")
	if err := repo.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	_ = context.Background() // ctx not used directly, db handles context
	db := repo.GetDB()

	// Get all mods
	var mods []*models.Mod
	if err := db.Find(&mods).Error; err != nil {
		return fmt.Errorf("failed to get mods: %w", err)
	}

	fmt.Printf("Found %d mods to migrate\n\n", len(mods))

	// Version info to fix
	versionFixes := map[string]struct {
		Version   string
		MCVersion string
	}{
		"bloodmagic":  {"3.0.0", "1.16.3"},
		"create":      {"0.5.1.f", "1.20.1"},
		"botania":     {"441", "1.20.1"},
		"mekanism":    {"10.4.16", "1.20.1"},
		"tconstruct":  {"3.10.2.92", "1.20.1"},
		"ars_nouveau": {"4.10.0", "1.20.1"},
	}

	totalSources := 0
	totalTranslations := 0

	for _, mod := range mods {
		fmt.Printf("Migrating %s...\n", mod.ID)

		// Get versions for this mod
		var versions []*models.ModVersion
		if err := db.Where("mod_id = ?", mod.ID).Find(&versions).Error; err != nil {
			return fmt.Errorf("failed to get versions for %s: %w", mod.ID, err)
		}

		if len(versions) == 0 {
			fmt.Printf("  No versions found, skipping\n")
			continue
		}

		// Fix version info if needed
		if fix, ok := versionFixes[mod.ID]; ok {
			for _, v := range versions {
				needsUpdate := false
				if v.Version == "${file.jarVersion}" || v.Version == "" {
					v.Version = fix.Version
					needsUpdate = true
				}
				if v.MCVersion == "" {
					v.MCVersion = fix.MCVersion
					needsUpdate = true
				}
				if !v.IsDefault {
					v.IsDefault = true
					needsUpdate = true
				}
				if needsUpdate && !*dryRun {
					if err := db.Save(v).Error; err != nil {
						return fmt.Errorf("failed to update version: %w", err)
					}
				}
				fmt.Printf("  Fixed version: %s (MC %s)\n", v.Version, v.MCVersion)
			}
		}

		// Use the first (or only) version
		version := versions[0]

		// Get existing translations for this version (using legacy query)
		var oldTranslations []struct {
			ID         int64
			Key        string
			SourceText string
			SourceLang string
			TargetText *string
			TargetLang string
			Status     string
			Translator *string
			Notes      *string
		}

		err := db.Table("translations").
			Select("id, key, source_text, source_lang, target_text, target_lang, status, translator, notes").
			Where("mod_version_id = ?", version.ID).
			Find(&oldTranslations).Error
		if err != nil {
			return fmt.Errorf("failed to get translations: %w", err)
		}

		fmt.Printf("  Found %d translations to migrate\n", len(oldTranslations))

		if *dryRun {
			totalSources += len(oldTranslations)
			totalTranslations += len(oldTranslations)
			continue
		}

		// Create sources and update translations
		for _, old := range oldTranslations {
			// Check if source already exists
			var existingSource models.TranslationSource
			err := db.Where("mod_id = ? AND key = ? AND source_text = ?", mod.ID, old.Key, old.SourceText).
				First(&existingSource).Error

			var sourceID int64
			if err != nil {
				// Create new source
				sourceLang := old.SourceLang
				if sourceLang == "" {
					sourceLang = "en_us"
				}
				source := &models.TranslationSource{
					ModID:      mod.ID,
					Key:        old.Key,
					SourceText: old.SourceText,
					SourceLang: sourceLang,
					IsCurrent:  true,
				}
				if err := db.Create(source).Error; err != nil {
					return fmt.Errorf("failed to create source for %s: %w", old.Key, err)
				}
				sourceID = source.ID
				totalSources++
			} else {
				sourceID = existingSource.ID
			}

			// Link source to version
			sv := &models.SourceVersion{
				SourceID:     sourceID,
				ModVersionID: version.ID,
			}
			db.Where("source_id = ? AND mod_version_id = ?", sourceID, version.ID).
				FirstOrCreate(sv)

			// Update translation to use source_id
			if err := db.Exec("UPDATE translations SET source_id = ? WHERE id = ?", sourceID, old.ID).Error; err != nil {
				return fmt.Errorf("failed to update translation %d: %w", old.ID, err)
			}
			totalTranslations++
		}

		fmt.Printf("  Migrated successfully\n")
	}

	if *dryRun {
		fmt.Printf("\n[DRY RUN] Would create %d sources and update %d translations\n", totalSources, totalTranslations)
	} else {
		fmt.Printf("\nMigration complete!\n")
		fmt.Printf("Created %d sources\n", totalSources)
		fmt.Printf("Updated %d translations\n", totalTranslations)
	}

	return nil
}
