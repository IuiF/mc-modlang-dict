package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// termFile represents the structure of a YAML term file.
type termFile struct {
	Scope string `yaml:"scope"`
	Terms []struct {
		Source   string   `yaml:"source"`
		Target   string   `yaml:"target"`
		Tags     []string `yaml:"tags,omitempty"`
		Priority int      `yaml:"priority,omitempty"`
		Context  string   `yaml:"context,omitempty"`
	} `yaml:"terms"`
}

// patternFile represents the structure of a YAML pattern file.
type patternFile struct {
	Patterns []struct {
		Pattern     string `yaml:"pattern"`
		Type        string `yaml:"type"`
		Parser      string `yaml:"parser"`
		Priority    int    `yaml:"priority,omitempty"`
		Required    bool   `yaml:"required,omitempty"`
		Description string `yaml:"description,omitempty"`
	} `yaml:"patterns"`
}

func runBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)

	var (
		dbPath   = fs.String("db", "moddict.db", "Database file path")
		dataDir  = fs.String("data", "data", "Data directory containing YAML files")
		termDir  = fs.String("terms", "", "Terms directory (default: data/terms)")
		patDir   = fs.String("patterns", "", "Patterns directory (default: data/patterns)")
		clean    = fs.Bool("clean", false, "Clean existing data before import")
	)

	fs.Usage = func() {
		fmt.Print(`Usage: moddict build [options]

Build translation database from YAML definition files.

Options:
`)
		fs.PrintDefaults()
		fmt.Print(`
Examples:
  moddict build
  moddict build -data ./my-data -db translations.db
  moddict build -clean
`)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Set default directories
	if *termDir == "" {
		*termDir = filepath.Join(*dataDir, "terms")
	}
	if *patDir == "" {
		*patDir = filepath.Join(*dataDir, "patterns")
	}

	// Clean database if requested
	if *clean {
		os.Remove(*dbPath)
		fmt.Println("Cleaned existing database")
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

	// Import terms
	termCount, err := importTerms(ctx, repo, *termDir)
	if err != nil {
		return fmt.Errorf("failed to import terms: %w", err)
	}
	fmt.Printf("Imported %d terms\n", termCount)

	// Import patterns
	patternCount, err := importPatterns(ctx, repo, *patDir)
	if err != nil {
		return fmt.Errorf("failed to import patterns: %w", err)
	}
	fmt.Printf("Imported %d patterns\n", patternCount)

	fmt.Printf("\nBuild complete: %s\n", *dbPath)
	return nil
}

func importTerms(ctx context.Context, repo *database.Repository, dir string) (int, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Terms directory not found: %s (skipping)\n", dir)
		return 0, nil
	}

	var totalCount int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var file termFile
		if err := yaml.Unmarshal(content, &file); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		terms := make([]*models.Term, 0, len(file.Terms))
		for _, t := range file.Terms {
			priority := t.Priority
			if priority == 0 {
				priority = 100
			}

			term := &models.Term{
				Scope:      file.Scope,
				SourceText: t.Source,
				TargetText: t.Target,
				SourceLang: "en_us",
				TargetLang: "ja_jp",
				Tags:       t.Tags,
				Priority:   priority,
			}

			if t.Context != "" {
				term.Context = &t.Context
			}

			terms = append(terms, term)
		}

		if len(terms) > 0 {
			if err := repo.BulkSaveTerms(ctx, terms); err != nil {
				return fmt.Errorf("failed to save terms from %s: %w", path, err)
			}
			totalCount += len(terms)
			fmt.Printf("  Loaded %d terms from %s\n", len(terms), filepath.Base(path))
		}

		return nil
	})

	return totalCount, err
}

func importPatterns(ctx context.Context, repo *database.Repository, dir string) (int, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Patterns directory not found: %s (skipping)\n", dir)
		return 0, nil
	}

	var totalCount int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		var file patternFile
		if err := yaml.Unmarshal(content, &file); err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		for _, p := range file.Patterns {
			priority := p.Priority
			if priority == 0 {
				priority = 100
			}

			pattern := &models.FilePattern{
				Scope:    "global",
				Pattern:  p.Pattern,
				Type:     p.Type,
				Parser:   p.Parser,
				Priority: priority,
				Required: p.Required,
			}

			if p.Description != "" {
				pattern.Description = &p.Description
			}

			if err := repo.SavePattern(ctx, pattern); err != nil {
				return fmt.Errorf("failed to save pattern: %w", err)
			}
			totalCount++
		}

		if len(file.Patterns) > 0 {
			fmt.Printf("  Loaded %d patterns from %s\n", len(file.Patterns), filepath.Base(path))
		}

		return nil
	})

	return totalCount, err
}
