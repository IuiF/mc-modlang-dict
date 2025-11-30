package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runFixSchema(args []string) error {
	fs := flag.NewFlagSet("fix-schema", flag.ExitOnError)

	var (
		dbPath = fs.String("db", "moddict.db", "Database file path")
		dryRun = fs.Bool("dry-run", false, "Show what would be done without making changes")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict fix-schema [options]

Fix translations to properly use the new source-based schema.
- Consolidates duplicate translations for same key+source_text
- Sets source_id for translations missing it
- Updates source_versions links

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

	_ = context.Background() // ctx not directly used
	db := repo.GetDB()

	// Get all mods
	var mods []*models.Mod
	if err := db.Find(&mods).Error; err != nil {
		return fmt.Errorf("failed to get mods: %w", err)
	}

	fmt.Printf("Fixing schema for %d mods...\n\n", len(mods))

	totalFixed := 0
	totalRemoved := 0

	for _, mod := range mods {
		fmt.Printf("Processing %s...\n", mod.ID)

		// Get all versions for this mod
		var versions []*models.ModVersion
		if err := db.Where("mod_id = ?", mod.ID).Find(&versions).Error; err != nil {
			return fmt.Errorf("failed to get versions: %w", err)
		}

		if len(versions) <= 1 {
			fmt.Printf("  Only %d version, skipping\n", len(versions))
			continue
		}

		// Build map of existing sources: key+source_text -> Source
		var sources []*models.TranslationSource
		if err := db.Where("mod_id = ?", mod.ID).Find(&sources).Error; err != nil {
			return fmt.Errorf("failed to get sources: %w", err)
		}

		sourceMap := make(map[string]*models.TranslationSource)
		for _, s := range sources {
			key := s.Key + "|" + s.SourceText
			sourceMap[key] = s
		}
		fmt.Printf("  Found %d existing sources\n", len(sources))

		// Get all translations for this mod
		var translations []*models.Translation
		if err := db.Joins("JOIN mod_versions ON translations.mod_version_id = mod_versions.id").
			Where("mod_versions.mod_id = ?", mod.ID).
			Find(&translations).Error; err != nil {
			return fmt.Errorf("failed to get translations: %w", err)
		}

		fmt.Printf("  Found %d translations\n", len(translations))

		// Group translations by key+source_text
		type TransGroup struct {
			Key        string
			SourceText string
			Items      []*models.Translation
		}
		groups := make(map[string]*TransGroup)

		for _, t := range translations {
			key := t.Key + "|" + t.SourceText
			if groups[key] == nil {
				groups[key] = &TransGroup{
					Key:        t.Key,
					SourceText: t.SourceText,
				}
			}
			groups[key].Items = append(groups[key].Items, t)
		}

		// Process each group
		for groupKey, group := range groups {
			if len(group.Items) <= 1 && group.Items[0].SourceID > 0 {
				// Already correct
				continue
			}

			// Find or create source
			source := sourceMap[groupKey]
			if source == nil {
				// Create new source
				source = &models.TranslationSource{
					ModID:      mod.ID,
					Key:        group.Key,
					SourceText: group.SourceText,
					SourceLang: "en_us",
					IsCurrent:  true,
				}
				if !*dryRun {
					if err := db.Create(source).Error; err != nil {
						fmt.Printf("  Warning: failed to create source for %s: %v\n", group.Key, err)
						continue
					}
				}
				sourceMap[groupKey] = source
			}

			// Find the "best" translation (prefer official, then translated with target)
			var bestTrans *models.Translation
			for _, t := range group.Items {
				if bestTrans == nil {
					bestTrans = t
					continue
				}
				// Prefer official
				if t.Translator != nil && *t.Translator == "official" {
					bestTrans = t
					break
				}
				// Prefer one with target text
				if t.TargetText != nil && *t.TargetText != "" && (bestTrans.TargetText == nil || *bestTrans.TargetText == "") {
					bestTrans = t
				}
			}

			// Update best translation to use source_id
			if bestTrans.SourceID != source.ID {
				if !*dryRun {
					bestTrans.SourceID = source.ID
					if err := db.Save(bestTrans).Error; err != nil {
						fmt.Printf("  Warning: failed to update translation %d: %v\n", bestTrans.ID, err)
						continue
					}
				}
				totalFixed++
			}

			// Link source to all versions
			for _, t := range group.Items {
				if !*dryRun {
					sv := &models.SourceVersion{
						SourceID:     source.ID,
						ModVersionID: t.ModVersionID,
					}
					db.Where("source_id = ? AND mod_version_id = ?", source.ID, t.ModVersionID).
						FirstOrCreate(sv)
				}
			}

			// Remove duplicate translations (keep only bestTrans)
			for _, t := range group.Items {
				if t.ID != bestTrans.ID {
					if !*dryRun {
						if err := db.Delete(t).Error; err != nil {
							fmt.Printf("  Warning: failed to delete duplicate translation %d: %v\n", t.ID, err)
							continue
						}
					}
					totalRemoved++
				}
			}
		}

		fmt.Printf("  Done\n")
	}

	if *dryRun {
		fmt.Printf("\n[DRY RUN] Would fix %d translations and remove %d duplicates\n", totalFixed, totalRemoved)
	} else {
		fmt.Printf("\nFixed %d translations and removed %d duplicates\n", totalFixed, totalRemoved)
	}

	return nil
}
