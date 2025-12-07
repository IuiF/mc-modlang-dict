package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
		modID      = fs.String("mod", "", "Mod ID to export (required)")
		targetLang = fs.String("lang", "ja_jp", "Target language code")
		format     = fs.String("format", "json", "Output format (json, merged, csv)")
		original   = fs.String("original", "", "Original lang file for merged export")
		status     = fs.String("status", "", "Filter by status (pending, translated, verified)")
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
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *modID == "" {
		fs.Usage()
		return fmt.Errorf("mod ID is required")
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
