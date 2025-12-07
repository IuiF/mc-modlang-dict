package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
	"gopkg.in/yaml.v3"
)

func runTranslate(args []string) error {
	fs := flag.NewFlagSet("translate", flag.ExitOnError)

	var (
		dbPath     = fs.String("db", "moddict.db", "Database file path")
		modID      = fs.String("mod", "", "Mod ID (required)")
		fromYAML   = fs.String("yaml", "", "Import translations from YAML file")
		fromJSON   = fs.String("json", "", "Import translations from JSON file (ja_jp.json format)")
		official   = fs.String("official", "", "Import official translations from ja_jp.json (translator=official, status=verified)")
		status     = fs.Bool("status", false, "Show translation status")
		pending    = fs.Bool("pending", false, "List pending translations")
		limit      = fs.Int("limit", 20, "Limit for listing")
		offset     = fs.Int("offset", 0, "Offset for listing (pagination)")
		exportJSON = fs.String("export", "", "Export pending to JSON file")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict translate [options]

Add or update translations in the database.
Translations are managed per-mod (not per-version) using the source_id schema.

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  # Show translation status for a mod
  moddict translate -mod bloodmagic -status

  # List pending translations
  moddict translate -mod bloodmagic -pending -limit 50

  # Export pending translations to JSON
  moddict translate -mod bloodmagic -export pending.json -limit 100

  # Import translations from ja_jp.json
  moddict translate -mod bloodmagic -json ./translations_ja.json

  # Import official translations (status=verified)
  moddict translate -mod bloodmagic -official ./lang/ja_jp.json

  # Import from YAML
  moddict translate -mod bloodmagic -yaml ./translations.yaml
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *modID == "" {
		fs.Usage()
		return fmt.Errorf("mod ID is required")
	}

	repo, err := database.NewRepository(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Verify mod exists
	if _, err := repo.GetMod(ctx, *modID); err != nil {
		return fmt.Errorf("mod %s not found: %w", *modID, err)
	}

	// Status mode
	if *status {
		return showStatus(ctx, repo, *modID)
	}

	// Export pending to JSON
	if *exportJSON != "" {
		return exportPendingJSON(ctx, repo, *modID, *exportJSON, *offset, *limit)
	}

	// Pending mode
	if *pending {
		return listPending(ctx, repo, *modID, *offset, *limit)
	}

	// Import from JSON
	if *fromJSON != "" {
		return importFromJSON(ctx, repo, *modID, *fromJSON)
	}

	// Import from YAML
	if *fromYAML != "" {
		return importFromYAML(ctx, repo, *modID, *fromYAML)
	}

	// Import official translations
	if *official != "" {
		return importOfficialJSON(ctx, repo, *modID, *official)
	}

	fs.Usage()
	return nil
}

func showStatus(ctx context.Context, repo *database.Repository, modID string) error {
	counts, err := repo.CountTranslationsByMod(ctx, modID)
	if err != nil {
		return err
	}

	pending := counts[models.StatusPending]
	translated := counts[models.StatusTranslated]
	verified := counts[models.StatusVerified]
	inherited := counts[models.StatusInherited]
	needsReview := counts[models.StatusNeedsReview]

	total := pending + translated + verified + inherited + needsReview
	done := translated + verified + inherited

	var progress float64
	if total > 0 {
		progress = float64(done) / float64(total) * 100
	}

	fmt.Printf("Translation Status for %s\n", modID)
	fmt.Printf("================================\n")
	fmt.Printf("Total keys:    %d\n", total)
	fmt.Printf("Pending:       %d\n", pending)
	fmt.Printf("Translated:    %d\n", translated)
	fmt.Printf("Inherited:     %d\n", inherited)
	fmt.Printf("Needs review:  %d\n", needsReview)
	fmt.Printf("Verified:      %d\n", verified)
	fmt.Printf("Progress:      %.1f%% (%d/%d)\n", progress, done, total)

	return nil
}

func listPending(ctx context.Context, repo *database.Repository, modID string, offset, limit int) error {
	filter := interfaces.TranslationFilter{
		Status: models.StatusPending,
		Offset: offset,
		Limit:  limit,
	}

	translations, err := repo.ListTranslationsWithSourceByMod(ctx, modID, filter)
	if err != nil {
		return err
	}

	fmt.Printf("Pending translations for %s (offset=%d, %d shown):\n\n", modID, offset, len(translations))
	for _, t := range translations {
		text := t.SourceText
		if len(text) > 60 {
			text = text[:60] + "..."
		}
		fmt.Printf("Key: %s\n", t.Key)
		fmt.Printf("Source: %s\n\n", text)
	}

	return nil
}

func exportPendingJSON(ctx context.Context, repo *database.Repository, modID, outPath string, offset, limit int) error {
	filter := interfaces.TranslationFilter{
		Status: models.StatusPending,
		Offset: offset,
		Limit:  limit,
	}

	translations, err := repo.ListTranslationsWithSourceByMod(ctx, modID, filter)
	if err != nil {
		return err
	}

	// Build JSON map
	result := make(map[string]string)
	for _, t := range translations {
		result[t.Key] = t.SourceText
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outPath, err)
	}

	fmt.Printf("Exported %d pending entries to %s (offset=%d)\n", len(translations), outPath, offset)
	return nil
}

func importFromJSON(ctx context.Context, repo *database.Repository, modID, jsonPath string) error {
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", jsonPath, err)
	}

	var translations map[string]string
	if err := json.Unmarshal(content, &translations); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Update translations by looking up source by mod_id and key
	var updated, notFound, skippedEmpty, skippedSameAsSource int
	for key, target := range translations {
		// Skip empty translations to prevent data corruption
		if target == "" {
			skippedEmpty++
			continue
		}

		source, err := repo.GetSourceByModAndKey(ctx, modID, key)
		if err != nil {
			fmt.Printf("Warning: failed to get source for %s: %v\n", key, err)
			continue
		}
		if source == nil {
			notFound++
			continue
		}

		// Skip if target is same as source (untranslated)
		if target == source.SourceText {
			skippedSameAsSource++
			continue
		}

		// Get translation by source_id
		trans, err := repo.GetTranslationBySourceID(ctx, source.ID)
		if err != nil {
			fmt.Printf("Warning: translation not found for %s: %v\n", key, err)
			continue
		}

		trans.TargetText = &target
		trans.Status = models.StatusTranslated
		if err := repo.SaveTranslation(ctx, trans); err != nil {
			fmt.Printf("Warning: failed to update %s: %v\n", key, err)
			continue
		}
		updated++
	}

	fmt.Printf("Updated %d translations from %s\n", updated, jsonPath)
	if skippedEmpty > 0 {
		fmt.Printf("Skipped empty translations: %d\n", skippedEmpty)
	}
	if skippedSameAsSource > 0 {
		fmt.Printf("Skipped same-as-source (untranslated): %d\n", skippedSameAsSource)
	}
	if notFound > 0 {
		fmt.Printf("Keys not found in DB: %d\n", notFound)
	}
	return nil
}

func importFromYAML(ctx context.Context, repo *database.Repository, modID, yamlPath string) error {
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", yamlPath, err)
	}

	// Parse YAML - support multiple formats
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract translations from YAML structure
	translations := extractTranslationsFromYAML(data)

	var updated, notFound, skippedEmpty, skippedSameAsSource int
	for key, target := range translations {
		// Skip empty translations to prevent data corruption
		if target == "" {
			skippedEmpty++
			continue
		}

		source, err := repo.GetSourceByModAndKey(ctx, modID, key)
		if err != nil {
			fmt.Printf("Warning: failed to get source for %s: %v\n", key, err)
			continue
		}
		if source == nil {
			notFound++
			continue
		}

		// Skip if target is same as source (untranslated)
		if target == source.SourceText {
			skippedSameAsSource++
			continue
		}

		trans, err := repo.GetTranslationBySourceID(ctx, source.ID)
		if err != nil {
			fmt.Printf("Warning: translation not found for %s: %v\n", key, err)
			continue
		}

		trans.TargetText = &target
		trans.Status = models.StatusTranslated
		if err := repo.SaveTranslation(ctx, trans); err != nil {
			fmt.Printf("Warning: failed to update %s: %v\n", key, err)
			continue
		}
		updated++
	}

	fmt.Printf("Updated %d translations from %s\n", updated, yamlPath)
	if skippedEmpty > 0 {
		fmt.Printf("Skipped empty translations: %d\n", skippedEmpty)
	}
	if skippedSameAsSource > 0 {
		fmt.Printf("Skipped same-as-source (untranslated): %d\n", skippedSameAsSource)
	}
	if notFound > 0 {
		fmt.Printf("Keys not found in DB: %d\n", notFound)
	}
	return nil
}

func extractTranslationsFromYAML(data map[string]interface{}) map[string]string {
	result := make(map[string]string)

	// Handle different YAML structures
	// Format 1: direct key-value
	// Format 2: translations array with key/target fields
	// Format 3: nested categories

	if translations, ok := data["translations"].([]interface{}); ok {
		for _, item := range translations {
			if m, ok := item.(map[string]interface{}); ok {
				key, _ := m["key"].(string)
				target, _ := m["target"].(string)
				if key != "" && target != "" {
					result[key] = target
				}
			}
		}
	}

	// Check for entries array (patchouli format)
	if entries, ok := data["entries"].([]interface{}); ok {
		for _, entry := range entries {
			if m, ok := entry.(map[string]interface{}); ok {
				extractPatchouliYAML(m, result)
			}
		}
	}

	return result
}

func extractPatchouliYAML(entry map[string]interface{}, result map[string]string) {
	id, _ := entry["id"].(string)
	if id == "" {
		return
	}

	keyPrefix := "patchouli:" + id

	// Name
	if name, ok := entry["name_target"].(string); ok && name != "" {
		result[keyPrefix+".name"] = name
	}

	// Pages
	if pages, ok := entry["pages"].([]interface{}); ok {
		for i, page := range pages {
			if m, ok := page.(map[string]interface{}); ok {
				if title, ok := m["title_target"].(string); ok && title != "" {
					result[fmt.Sprintf("%s.page%d.title", keyPrefix, i)] = title
				}
				if text, ok := m["text_target"].(string); ok && text != "" {
					result[fmt.Sprintf("%s.page%d.text", keyPrefix, i)] = text
				}
				// Also support "target" field
				if target, ok := m["target"].(string); ok && target != "" {
					result[fmt.Sprintf("%s.page%d.text", keyPrefix, i)] = target
				}
			}
		}
	}
}

// importOfficialJSON imports official translations from ja_jp.json
// Sets translator=official, status=verified for imported translations
// Skips entries where target_text equals source_text (untranslated in official file)
func importOfficialJSON(ctx context.Context, repo *database.Repository, modID, jsonPath string) error {
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", jsonPath, err)
	}

	var officialTranslations map[string]string
	if err := json.Unmarshal(content, &officialTranslations); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Track statistics
	var updated, alreadyOfficial, notFound, skippedEmpty, skippedSameAsSource int
	officialStr := "official"

	for key, target := range officialTranslations {
		// Skip empty translations
		if target == "" {
			skippedEmpty++
			continue
		}

		source, err := repo.GetSourceByModAndKey(ctx, modID, key)
		if err != nil {
			fmt.Printf("Warning: failed to get source for %s: %v\n", key, err)
			continue
		}
		if source == nil {
			notFound++
			continue
		}

		// Skip if target is same as source (untranslated in official file)
		if target == source.SourceText {
			skippedSameAsSource++
			continue
		}

		trans, err := repo.GetTranslationBySourceID(ctx, source.ID)
		if err != nil {
			fmt.Printf("Warning: translation not found for %s: %v\n", key, err)
			continue
		}

		// Check if already marked as official
		if trans.Translator != nil && *trans.Translator == "official" {
			alreadyOfficial++
			continue
		}

		// Update with official translation
		trans.TargetText = &target
		trans.Status = models.StatusVerified
		trans.Translator = &officialStr
		if err := repo.SaveTranslation(ctx, trans); err != nil {
			fmt.Printf("Warning: failed to update %s: %v\n", key, err)
			continue
		}
		updated++
	}

	// Print summary
	fmt.Printf("Official translation import from %s\n", jsonPath)
	fmt.Printf("========================================\n")
	fmt.Printf("Total official keys:    %d\n", len(officialTranslations))
	fmt.Printf("Updated to official:    %d\n", updated)
	fmt.Printf("Already official:       %d\n", alreadyOfficial)
	fmt.Printf("Skipped empty:          %d\n", skippedEmpty)
	fmt.Printf("Skipped same-as-source: %d\n", skippedSameAsSource)
	fmt.Printf("Keys not in DB:         %d\n", notFound)

	return nil
}
