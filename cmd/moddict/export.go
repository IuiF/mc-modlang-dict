package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/internal/export"
	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ExitOnError)

	var (
		dbPath     = fs.String("db", "moddict.db", "Database file path")
		outputDir  = fs.String("out", "workspace/exports", "Output directory")
		modID      = fs.String("mod", "", "Mod ID to export")
		targetLang = fs.String("lang", "ja_jp", "Target language code")
		format     = fs.String("format", "json", "Output format (json, merged, csv, resourcepack)")
		original   = fs.String("original", "", "Original lang file for merged export")
		status     = fs.String("status", "", "Filter by status (pending, translated, verified)")
		all        = fs.Bool("all", false, "Export all mods to a single combined CSV file")
		perMod     = fs.Bool("per-mod", false, "Export each mod to separate CSV files (use with -all)")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict export [options]

Export translations to various formats.
Translations are exported per-mod using the source_id schema.

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  moddict export -mod create -out output/
  moddict export -mod botania -format merged -original en_us.json
  moddict export -mod create -status translated
  moddict export -all -out translations/       # Export all mods to combined CSV
  moddict export -all -per-mod -out data/translations/  # Export each mod to separate CSV
  moddict export -all -format resourcepack -out my_pack/  # Export as Minecraft resource pack
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Handle -all flag
	if *all {
		if *format == "resourcepack" {
			return runExportResourcePack(*dbPath, *outputDir, *targetLang, *status)
		}
		return runExportAllCSV(*dbPath, *outputDir, *targetLang, *status, *perMod)
	}

	if *modID == "" {
		fs.Usage()
		return fmt.Errorf("mod ID is required (or use -all for combined export)")
	}

	// Open database
	repo, err := database.NewRepository(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Get mod info
	mod, err := repo.GetMod(ctx, *modID)
	if err != nil {
		return fmt.Errorf("mod not found: %s", *modID)
	}

	// Get translations using new source-based query
	filter := interfaces.TranslationFilter{
		TargetLang: *targetLang,
	}
	if *status != "" {
		filter.Status = *status
	}

	translations, err := repo.ListTranslationsWithSourceByMod(ctx, *modID, filter)
	if err != nil {
		return fmt.Errorf("failed to get translations: %w", err)
	}

	if len(translations) == 0 {
		fmt.Println("No translations found matching criteria")
		return nil
	}

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	exporter := export.NewExporter()
	outputPath := filepath.Join(*outputDir, fmt.Sprintf("%s_%s.json", *modID, *targetLang))

	switch *format {
	case "json":
		fmt.Printf("Exporting %d translations to %s...\n", len(translations), outputPath)
		if err := exporter.ExportJSON(translations, outputPath); err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

	case "merged":
		if *original == "" {
			return fmt.Errorf("-original flag is required for merged export")
		}

		originalContent, err := os.ReadFile(*original)
		if err != nil {
			return fmt.Errorf("failed to read original file: %w", err)
		}

		fmt.Printf("Exporting merged translations to %s...\n", outputPath)
		if err := exporter.ExportMerged(originalContent, translations, outputPath); err != nil {
			return fmt.Errorf("failed to export: %w", err)
		}

	case "csv":
		csvPath := filepath.Join(*outputDir, fmt.Sprintf("%s_%s.csv", *modID, *targetLang))
		fmt.Printf("Exporting %d translations to %s...\n", len(translations), csvPath)
		if err := exportCSV(translations, csvPath); err != nil {
			return fmt.Errorf("failed to export CSV: %w", err)
		}
		outputPath = csvPath

	default:
		return fmt.Errorf("unknown format: %s", *format)
	}

	// Print summary
	var translatedCount int
	for _, t := range translations {
		if t.TargetText != nil {
			translatedCount++
		}
	}

	fmt.Printf("\nExport complete for %s (%s)\n", mod.DisplayName, mod.ID)
	fmt.Printf("Total entries: %d\n", len(translations))
	fmt.Printf("Translated: %d\n", translatedCount)
	fmt.Printf("Output: %s\n", outputPath)

	return nil
}

// runExportAllCSV exports all mods' translations to CSV files
// If perMod is true, exports each mod to a separate CSV file
// Otherwise, exports all mods to a single combined CSV file
func runExportAllCSV(dbPath, outputDir, targetLang, status string, perMod bool) error {
	// Open database
	repo, err := database.NewRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Get all mods
	mods, err := repo.ListMods(ctx, interfaces.ModFilter{})
	if err != nil {
		return fmt.Errorf("failed to list mods: %w", err)
	}

	if len(mods) == 0 {
		fmt.Println("No mods found in database")
		return nil
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if perMod {
		return runExportPerModCSV(repo, ctx, mods, outputDir, targetLang, status)
	}

	return runExportCombinedCSV(repo, ctx, mods, outputDir, targetLang, status)
}

// runExportPerModCSV exports each mod to a separate CSV file
func runExportPerModCSV(repo *database.Repository, ctx context.Context, mods []*models.Mod, outputDir, targetLang, status string) error {
	totalMods := 0
	totalRows := 0
	translatedCount := 0

	for _, mod := range mods {
		filter := interfaces.TranslationFilter{
			TargetLang: targetLang,
		}
		if status != "" {
			filter.Status = status
		}

		translations, err := repo.ListTranslationsWithSourceByMod(ctx, mod.ID, filter)
		if err != nil {
			fmt.Printf("Warning: failed to get translations for %s: %v\n", mod.ID, err)
			continue
		}

		if len(translations) == 0 {
			continue
		}

		// Create CSV file for this mod
		csvPath := filepath.Join(outputDir, fmt.Sprintf("%s.csv", mod.ID))
		if err := exportCSV(translations, csvPath); err != nil {
			fmt.Printf("Warning: failed to export %s: %v\n", mod.ID, err)
			continue
		}

		totalMods++
		totalRows += len(translations)
		for _, t := range translations {
			if t.TargetText != nil {
				translatedCount++
			}
		}
	}

	fmt.Printf("\nPer-Mod CSV Export Complete\n")
	fmt.Printf("Mods exported: %d\n", totalMods)
	fmt.Printf("Total entries: %d\n", totalRows)
	fmt.Printf("Translated: %d\n", translatedCount)
	fmt.Printf("Output directory: %s\n", outputDir)

	return nil
}

// runExportCombinedCSV exports all mods to a single combined CSV file
func runExportCombinedCSV(repo *database.Repository, ctx context.Context, mods []*models.Mod, outputDir, targetLang, status string) error {
	// Create combined CSV file
	csvPath := filepath.Join(outputDir, "all_translations.csv")
	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write UTF-8 BOM for Excel compatibility
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header - only key, source_text, target_text
	header := []string{"key", "source_text", "target_text"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	totalRows := 0
	translatedCount := 0

	// Process each mod
	for _, mod := range mods {
		filter := interfaces.TranslationFilter{
			TargetLang: targetLang,
		}
		if status != "" {
			filter.Status = status
		}

		translations, err := repo.ListTranslationsWithSourceByMod(ctx, mod.ID, filter)
		if err != nil {
			fmt.Printf("Warning: failed to get translations for %s: %v\n", mod.ID, err)
			continue
		}

		// Write rows for this mod
		for _, t := range translations {
			targetText := ""
			if t.TargetText != nil {
				targetText = *t.TargetText
				translatedCount++
			}
			row := []string{t.Key, t.SourceText, targetText}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write row: %w", err)
			}
			totalRows++
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	fmt.Printf("\nCombined CSV Export Complete\n")
	fmt.Printf("Mods processed: %d\n", len(mods))
	fmt.Printf("Total entries: %d\n", totalRows)
	fmt.Printf("Translated: %d\n", translatedCount)
	fmt.Printf("Output: %s\n", csvPath)

	return nil
}

// exportCSV exports translations to a CSV file
// Format: key,source_text,target_text
func exportCSV(translations []*models.TranslationWithSource, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write UTF-8 BOM for Excel compatibility
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"key", "source_text", "target_text"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write rows
	for _, t := range translations {
		targetText := ""
		if t.TargetText != nil {
			targetText = *t.TargetText
		}
		row := []string{t.Key, t.SourceText, targetText}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// runExportResourcePack exports all translations as a Minecraft resource pack
// Structure:
//
//	<outputDir>/
//	├── pack.mcmeta
//	└── assets/
//	    └── <mod_id>/
//	        └── lang/
//	            └── <targetLang>.json
func runExportResourcePack(dbPath, outputDir, targetLang, status string) error {
	// Open database
	repo, err := database.NewRepository(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	ctx := context.Background()

	// Get all mods
	mods, err := repo.ListMods(ctx, interfaces.ModFilter{})
	if err != nil {
		return fmt.Errorf("failed to list mods: %w", err)
	}

	if len(mods) == 0 {
		fmt.Println("No mods found in database")
		return nil
	}

	// Create output directory structure
	assetsDir := filepath.Join(outputDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Create pack.mcmeta
	if err := createPackMcmeta(outputDir); err != nil {
		return fmt.Errorf("failed to create pack.mcmeta: %w", err)
	}

	totalMods := 0
	totalEntries := 0
	translatedCount := 0

	// Process each mod
	for _, mod := range mods {
		filter := interfaces.TranslationFilter{
			TargetLang: targetLang,
		}
		if status != "" {
			filter.Status = status
		}

		translations, err := repo.ListTranslationsWithSourceByMod(ctx, mod.ID, filter)
		if err != nil {
			fmt.Printf("Warning: failed to get translations for %s: %v\n", mod.ID, err)
			continue
		}

		// Only include mods with translated entries
		translatedEntries := make(map[string]string)
		for _, t := range translations {
			if t.TargetText != nil && *t.TargetText != "" {
				translatedEntries[t.Key] = *t.TargetText
				translatedCount++
			}
		}

		if len(translatedEntries) == 0 {
			continue
		}

		// Create mod lang directory
		langDir := filepath.Join(assetsDir, mod.ID, "lang")
		if err := os.MkdirAll(langDir, 0755); err != nil {
			fmt.Printf("Warning: failed to create directory for %s: %v\n", mod.ID, err)
			continue
		}

		// Write lang file
		langFile := filepath.Join(langDir, targetLang+".json")
		if err := writeLangJSON(translatedEntries, langFile); err != nil {
			fmt.Printf("Warning: failed to write lang file for %s: %v\n", mod.ID, err)
			continue
		}

		totalMods++
		totalEntries += len(translatedEntries)
	}

	fmt.Printf("\nResource Pack Export Complete\n")
	fmt.Printf("Mods exported: %d\n", totalMods)
	fmt.Printf("Total entries: %d\n", totalEntries)
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Printf("\nTo use: Copy '%s' to your Minecraft resourcepacks folder\n", outputDir)

	return nil
}

// createPackMcmeta creates the pack.mcmeta file for the resource pack
func createPackMcmeta(outputDir string) error {
	packMcmeta := map[string]interface{}{
		"pack": map[string]interface{}{
			"pack_format": 15, // 1.20.x format
			"description": "Mod翻訳リソースパック - Generated by moddict",
		},
	}

	data, err := json.MarshalIndent(packMcmeta, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, "pack.mcmeta"), data, 0644)
}

// writeLangJSON writes a language JSON file with sorted keys
func writeLangJSON(entries map[string]string, outputPath string) error {
	// Sort keys for consistent output
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build ordered map for JSON output
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write as formatted JSON
	file.WriteString("{\n")
	for i, key := range keys {
		// Escape special characters in value
		value := entries[key]
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		file.WriteString(fmt.Sprintf("  %q: %s", key, string(jsonValue)))
		if i < len(keys)-1 {
			file.WriteString(",")
		}
		file.WriteString("\n")
	}
	file.WriteString("}\n")

	return nil
}
