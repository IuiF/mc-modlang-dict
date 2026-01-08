package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/iuif/minecraft-mod-dictionary/internal/analyzer"
	"github.com/iuif/minecraft-mod-dictionary/internal/database"
)

func runAnalyze(args []string) error {
	// Extract subcommand first (before flags)
	subcommand := "all"
	flagArgs := args

	if len(args) > 0 {
		switch args[0] {
		case "consistency", "phrases", "terms", "all":
			subcommand = args[0]
			flagArgs = args[1:]
		case "-h", "--help", "-help":
			// Let the flag parser handle help
		default:
			// Check if first arg looks like a flag
			if len(args[0]) > 0 && args[0][0] != '-' {
				return fmt.Errorf("unknown subcommand: %s", args[0])
			}
		}
	}

	fs := flag.NewFlagSet("analyze", flag.ExitOnError)

	var (
		dbPath   = fs.String("db", "moddict.db", "Database file path")
		modID    = fs.String("mod", "", "Target mod ID (empty for all mods)")
		format   = fs.String("format", "summary", "Output format: summary, json, csv")
		outPath  = fs.String("out", "", "Output file path (empty for stdout)")
		minCount = fs.Int("min-count", 3, "Minimum occurrence count for phrase detection")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict analyze [subcommand] [options]

Analyze translation consistency and discover phrase patterns.

Subcommands:
  consistency  Check for same source text with different translations
  phrases      Discover and analyze phrase patterns (N-gram mining)
  terms        Check translations against term dictionary
  all          Run all analysis types (default)

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  moddict analyze consistency -mod create
  moddict analyze phrases -format json -out /tmp/phrases.json
  moddict analyze terms -mod botania -format csv
  moddict analyze all -mod mekanism
`)
	}

	if err := fs.Parse(flagArgs); err != nil {
		return err
	}

	repo, err := database.NewRepository(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer repo.Close()

	ctx := context.Background()

	opts := analyzer.AnalysisOptions{
		ModID:      *modID,
		MinCount:   *minCount,
		Format:     *format,
		OutputPath: *outPath,
	}

	a := analyzer.New(repo)

	var result *analyzer.AnalysisResult

	switch subcommand {
	case "consistency":
		result, err = a.AnalyzeConsistency(ctx, opts)
	case "phrases":
		result, err = a.AnalyzePhrases(ctx, opts)
	case "terms":
		result, err = a.AnalyzeTerms(ctx, opts)
	case "all":
		result, err = a.AnalyzeAll(ctx, opts)
	}

	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Output result
	var output []byte
	switch *format {
	case "json":
		output, err = analyzer.FormatJSON(result)
	case "csv":
		output, err = analyzer.FormatCSV(result)
	default:
		output, err = analyzer.FormatSummary(result)
	}

	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	if *outPath != "" {
		if err := os.WriteFile(*outPath, output, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Output written to: %s\n", *outPath)
	} else {
		fmt.Print(string(output))
	}

	return nil
}
