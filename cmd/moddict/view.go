package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
)

func runView(args []string) error {
	fs := flag.NewFlagSet("view", flag.ExitOnError)

	var (
		dbPath     = fs.String("db", "moddict.db", "Database file path")
		modID      = fs.String("mod", "", "Mod ID to filter (optional, shows all mods if not specified)")
		status     = fs.String("status", "", "Filter by status (pending, translated, verified, inherited, needs_review)")
		search     = fs.String("search", "", "Search in source_text or target_text")
		_          = fs.String("translator", "", "Filter by translator (e.g., 'lm:*' for all LM translations)")
		limit      = fs.Int("limit", 50, "Number of entries to show")
		offset     = fs.Int("offset", 0, "Offset for pagination")
		compact    = fs.Bool("compact", false, "Compact view (key and target only)")
		_          = fs.Bool("recent", false, "Order by most recently updated")
		id         = fs.Int64("id", 0, "Show specific translation by ID")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict view [options]

View translations in the database.

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  # View all translations (first 50)
  moddict view

  # View translations for a specific mod
  moddict view -mod bloodmagic

  # View translated entries only
  moddict view -status translated

  # Search for specific text
  moddict view -search "エネルギー"

  # Pagination
  moddict view -offset 100 -limit 20

  # View specific translation by ID
  moddict view -id 12345

  # Compact view
  moddict view -mod create -compact
`)
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

	// View specific ID
	if *id > 0 {
		return viewByID(ctx, repo, *id)
	}

	// View with filters
	return viewTranslations(ctx, repo, *modID, *status, *search, *offset, *limit, *compact)
}

func viewByID(ctx context.Context, repo *database.Repository, id int64) error {
	// Direct SQL query to get translation with source info
	var result struct {
		ID         int64
		SourceID   int64
		Key        string
		SourceText string
		TargetText *string
		Status     string
		Translator *string
		ModID      string
	}

	err := repo.DB().WithContext(ctx).
		Table("translations").
		Select(`
			translations.id,
			translations.source_id,
			translation_sources.key,
			translation_sources.source_text,
			translations.target_text,
			translations.status,
			translations.translator,
			translation_sources.mod_id
		`).
		Joins("JOIN translation_sources ON translation_sources.id = translations.source_id").
		Where("translations.id = ?", id).
		Scan(&result).Error

	if err != nil {
		return fmt.Errorf("failed to get translation: %w", err)
	}

	if result.ID == 0 {
		return fmt.Errorf("translation with ID %d not found", id)
	}

	fmt.Printf("Translation ID: %d\n", result.ID)
	fmt.Printf("=====================================\n")
	fmt.Printf("Mod:        %s\n", result.ModID)
	fmt.Printf("Key:        %s\n", result.Key)
	fmt.Printf("Status:     %s\n", result.Status)
	if result.Translator != nil {
		fmt.Printf("Translator: %s\n", *result.Translator)
	}
	fmt.Printf("\n--- Source Text ---\n%s\n", result.SourceText)
	if result.TargetText != nil {
		fmt.Printf("\n--- Target Text (ja_jp) ---\n%s\n", *result.TargetText)
	} else {
		fmt.Printf("\n--- Target Text (ja_jp) ---\n(not translated)\n")
	}

	return nil
}

func viewTranslations(ctx context.Context, repo *database.Repository, modID, status, search string, offset, limit int, compact bool) error {
	// Build query
	query := repo.DB().WithContext(ctx).
		Table("translations").
		Select(`
			translations.id,
			translation_sources.mod_id,
			translation_sources.key,
			translation_sources.source_text,
			translations.target_text,
			translations.status
		`).
		Joins("JOIN translation_sources ON translation_sources.id = translations.source_id")

	// Apply filters
	if modID != "" {
		query = query.Where("translation_sources.mod_id = ?", modID)
	}
	if status != "" {
		query = query.Where("translations.status = ?", status)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"translation_sources.source_text LIKE ? OR translations.target_text LIKE ?",
			searchPattern, searchPattern,
		)
	}

	// Get total count
	var totalCount int64
	countQuery := repo.DB().WithContext(ctx).
		Table("translations").
		Joins("JOIN translation_sources ON translation_sources.id = translations.source_id")

	if modID != "" {
		countQuery = countQuery.Where("translation_sources.mod_id = ?", modID)
	}
	if status != "" {
		countQuery = countQuery.Where("translations.status = ?", status)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		countQuery = countQuery.Where(
			"translation_sources.source_text LIKE ? OR translations.target_text LIKE ?",
			searchPattern, searchPattern,
		)
	}
	countQuery.Count(&totalCount)

	// Apply pagination
	query = query.Order("translations.id").Offset(offset).Limit(limit)

	// Execute query
	var results []struct {
		ID         int64
		ModID      string
		Key        string
		SourceText string
		TargetText *string
		Status     string
	}

	if err := query.Scan(&results).Error; err != nil {
		return fmt.Errorf("failed to query translations: %w", err)
	}

	// Print header
	filterInfo := []string{}
	if modID != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("mod=%s", modID))
	}
	if status != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("status=%s", status))
	}
	if search != "" {
		filterInfo = append(filterInfo, fmt.Sprintf("search=%q", search))
	}

	filterStr := ""
	if len(filterInfo) > 0 {
		filterStr = " (" + strings.Join(filterInfo, ", ") + ")"
	}

	fmt.Printf("Translations%s\n", filterStr)
	fmt.Printf("Showing %d-%d of %d entries\n", offset+1, offset+len(results), totalCount)
	fmt.Println(strings.Repeat("=", 80))

	if len(results) == 0 {
		fmt.Println("No translations found.")
		return nil
	}

	// Print results
	if compact {
		for _, r := range results {
			target := "(pending)"
			if r.TargetText != nil && *r.TargetText != "" {
				target = *r.TargetText
			}
			// Truncate if too long
			if len(target) > 60 {
				target = target[:57] + "..."
			}
			fmt.Printf("[%d] %s\n    → %s\n", r.ID, truncate(r.Key, 70), target)
		}
	} else {
		for i, r := range results {
			if i > 0 {
				fmt.Println(strings.Repeat("-", 80))
			}
			fmt.Printf("ID: %d | Mod: %s | Status: %s\n", r.ID, r.ModID, r.Status)
			fmt.Printf("Key: %s\n", r.Key)
			fmt.Printf("Source: %s\n", truncate(r.SourceText, 200))
			if r.TargetText != nil && *r.TargetText != "" {
				fmt.Printf("Target: %s\n", truncate(*r.TargetText, 200))
			} else {
				fmt.Printf("Target: (not translated)\n")
			}
		}
	}

	// Print pagination info
	fmt.Println(strings.Repeat("=", 80))
	if int64(offset+limit) < totalCount {
		fmt.Printf("Next: moddict view -offset %d -limit %d", offset+limit, limit)
		if modID != "" {
			fmt.Printf(" -mod %s", modID)
		}
		if status != "" {
			fmt.Printf(" -status %s", status)
		}
		if search != "" {
			fmt.Printf(" -search %q", search)
		}
		fmt.Println()
	}

	return nil
}

func truncate(s string, maxLen int) string {
	// Replace newlines with spaces for display
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
