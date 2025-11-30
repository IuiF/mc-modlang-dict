package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runRepair(args []string) error {
	fs := flag.NewFlagSet("repair", flag.ExitOnError)

	var (
		dbPath = fs.String("db", "moddict.db", "Database file path")
		dryRun = fs.Bool("dry-run", false, "Show what would be repaired without making changes")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict repair [options]

Repair database inconsistencies:
1. Set is_default=true for only the latest version of each mod
2. Merge duplicate sources (same mod_id + key + source_text)
3. Clean up orphaned data

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

	ctx := context.Background()
	db := repo.GetDB()

	fmt.Println("=== Database Repair ===")
	fmt.Println()

	// 1. Fix is_default for mod_versions
	fmt.Println("Step 1: Fixing mod_versions.is_default...")
	var mods []*models.Mod
	if err := db.Find(&mods).Error; err != nil {
		return fmt.Errorf("failed to get mods: %w", err)
	}

	versionsFixed := 0
	for _, mod := range mods {
		var versions []*models.ModVersion
		if err := db.Where("mod_id = ?", mod.ID).Order("id DESC").Find(&versions).Error; err != nil {
			return fmt.Errorf("failed to get versions for %s: %w", mod.ID, err)
		}

		if len(versions) == 0 {
			continue
		}

		// Count current defaults
		defaultCount := 0
		for _, v := range versions {
			if v.IsDefault {
				defaultCount++
			}
		}

		// Fix: only the latest (first in DESC order) should be default
		for i, v := range versions {
			shouldBeDefault := (i == 0)
			if v.IsDefault != shouldBeDefault {
				if *dryRun {
					fmt.Printf("  Would update %s v%s: is_default=%v -> %v\n", mod.ID, v.Version, v.IsDefault, shouldBeDefault)
				} else {
					if err := db.Model(v).Update("is_default", shouldBeDefault).Error; err != nil {
						return fmt.Errorf("failed to update version: %w", err)
					}
				}
				versionsFixed++
			}
		}
	}
	fmt.Printf("  Versions to fix: %d\n\n", versionsFixed)

	// 2. Find and merge duplicate sources
	fmt.Println("Step 2: Finding duplicate sources (same mod_id + key + source_text)...")
	type DuplicateGroup struct {
		ModID      string
		Key        string
		SourceText string
		Count      int64
		MinID      int64
	}
	var duplicates []DuplicateGroup
	err = db.Raw(`
		SELECT mod_id, key, source_text, COUNT(*) as count, MIN(id) as min_id
		FROM translation_sources
		GROUP BY mod_id, key, source_text
		HAVING count > 1
	`).Scan(&duplicates).Error
	if err != nil {
		return fmt.Errorf("failed to find duplicates: %w", err)
	}

	fmt.Printf("  Found %d duplicate groups\n", len(duplicates))

	sourcesMerged := 0
	for _, dup := range duplicates {
		if *dryRun {
			fmt.Printf("  Would merge %d duplicates for key: %s\n", dup.Count, dup.Key)
			sourcesMerged += int(dup.Count - 1)
			continue
		}

		// Get all source IDs for this duplicate group
		var sourceIDs []int64
		db.Raw(`
			SELECT id FROM translation_sources
			WHERE mod_id = ? AND key = ? AND source_text = ?
			ORDER BY id
		`, dup.ModID, dup.Key, dup.SourceText).Scan(&sourceIDs)

		if len(sourceIDs) <= 1 {
			continue
		}

		keepID := sourceIDs[0]
		deleteIDs := sourceIDs[1:]

		// Update source_versions to point to the kept source
		for _, delID := range deleteIDs {
			db.Exec(`
				UPDATE source_versions SET source_id = ?
				WHERE source_id = ? AND NOT EXISTS (
					SELECT 1 FROM source_versions sv2
					WHERE sv2.source_id = ? AND sv2.mod_version_id = source_versions.mod_version_id
				)
			`, keepID, delID, keepID)

			// Delete duplicate source_versions
			db.Exec(`DELETE FROM source_versions WHERE source_id = ?`, delID)

			// Update translations to point to the kept source (keep the first translation)
			var transCount int64
			db.Raw(`SELECT COUNT(*) FROM translations WHERE source_id = ?`, keepID).Scan(&transCount)
			if transCount == 0 {
				// No translation for kept source, update one from deleted
				db.Exec(`
					UPDATE translations SET source_id = ?
					WHERE source_id = ? AND id = (
						SELECT MIN(id) FROM translations WHERE source_id = ?
					)
				`, keepID, delID, delID)
			}
			// Delete remaining translations for deleted source
			db.Exec(`DELETE FROM translations WHERE source_id = ?`, delID)

			// Delete the duplicate source
			db.Exec(`DELETE FROM translation_sources WHERE id = ?`, delID)
			sourcesMerged++
		}
	}
	fmt.Printf("  Sources merged: %d\n\n", sourcesMerged)

	// 3. Clean up sources without translations
	fmt.Println("Step 3: Finding orphaned sources (no translation)...")
	var orphanedSources int64
	db.Raw(`
		SELECT COUNT(*) FROM translation_sources s
		WHERE NOT EXISTS (SELECT 1 FROM translations t WHERE t.source_id = s.id)
	`).Scan(&orphanedSources)
	fmt.Printf("  Orphaned sources: %d\n", orphanedSources)

	if orphanedSources > 0 && !*dryRun {
		// Create pending translations for orphaned sources
		db.Exec(`
			INSERT INTO translations (source_id, target_lang, status, created_at, updated_at)
			SELECT s.id, 'ja_jp', 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
			FROM translation_sources s
			WHERE NOT EXISTS (SELECT 1 FROM translations t WHERE t.source_id = s.id)
		`)
		fmt.Printf("  Created pending translations for orphaned sources\n")
	}

	// 4. Verify counts
	fmt.Println("\nStep 4: Verifying data integrity...")

	// Check for sources without version links
	var unlinkedSources int64
	db.Raw(`
		SELECT COUNT(*) FROM translation_sources s
		WHERE NOT EXISTS (SELECT 1 FROM source_versions sv WHERE sv.source_id = s.id)
	`).Scan(&unlinkedSources)
	fmt.Printf("  Sources without version links: %d\n", unlinkedSources)

	// Link unlinked sources to default version
	if unlinkedSources > 0 && !*dryRun {
		for _, mod := range mods {
			defaultVersion, err := repo.GetDefaultVersion(ctx, mod.ID)
			if err != nil {
				continue
			}
			db.Exec(`
				INSERT INTO source_versions (source_id, mod_version_id, created_at)
				SELECT s.id, ?, CURRENT_TIMESTAMP
				FROM translation_sources s
				WHERE s.mod_id = ?
				AND NOT EXISTS (SELECT 1 FROM source_versions sv WHERE sv.source_id = s.id)
			`, defaultVersion.ID, mod.ID)
		}
		fmt.Printf("  Linked orphaned sources to default versions\n")
	}

	// 5. Link all sources to default version (for sources linked to other versions but not default)
	fmt.Println("\nStep 5: Ensuring all sources are linked to default version...")
	var sourcesNotLinkedToDefault int64
	for _, mod := range mods {
		defaultVersion, err := repo.GetDefaultVersion(ctx, mod.ID)
		if err != nil {
			continue
		}

		var count int64
		db.Raw(`
			SELECT COUNT(*) FROM translation_sources s
			WHERE s.mod_id = ?
			AND NOT EXISTS (
				SELECT 1 FROM source_versions sv
				WHERE sv.source_id = s.id AND sv.mod_version_id = ?
			)
		`, mod.ID, defaultVersion.ID).Scan(&count)

		if count > 0 {
			fmt.Printf("  %s: %d sources not linked to default version\n", mod.ID, count)
			sourcesNotLinkedToDefault += count

			if !*dryRun {
				db.Exec(`
					INSERT INTO source_versions (source_id, mod_version_id, created_at)
					SELECT s.id, ?, CURRENT_TIMESTAMP
					FROM translation_sources s
					WHERE s.mod_id = ?
					AND NOT EXISTS (
						SELECT 1 FROM source_versions sv
						WHERE sv.source_id = s.id AND sv.mod_version_id = ?
					)
				`, defaultVersion.ID, mod.ID, defaultVersion.ID)
			}
		}
	}
	if sourcesNotLinkedToDefault == 0 {
		fmt.Printf("  All sources already linked to default versions\n")
	} else if !*dryRun {
		fmt.Printf("  Linked %d sources to default versions\n", sourcesNotLinkedToDefault)
	}

	// Final stats
	fmt.Println("\n=== Summary ===")
	var totalSources, totalTranslations int64
	db.Raw(`SELECT COUNT(*) FROM translation_sources`).Scan(&totalSources)
	db.Raw(`SELECT COUNT(*) FROM translations`).Scan(&totalTranslations)
	fmt.Printf("Total sources: %d\n", totalSources)
	fmt.Printf("Total translations: %d\n", totalTranslations)

	if *dryRun {
		fmt.Println("\n[DRY RUN] No changes made.")
	} else {
		fmt.Println("\nRepair complete!")
	}

	return nil
}
