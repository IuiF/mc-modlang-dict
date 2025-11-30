// Package export provides functionality for exporting translations to various formats.
package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// Exporter handles translation export to various formats.
type Exporter struct{}

// NewExporter creates a new exporter.
func NewExporter() *Exporter {
	return &Exporter{}
}

// ExportJSON exports translations to a JSON lang file.
// Uses TranslationWithSource which contains Key from translation_sources.
func (e *Exporter) ExportJSON(translations []*models.TranslationWithSource, destPath string) error {
	result := make(map[string]string)

	for _, trans := range translations {
		if trans.TargetText != nil && trans.Key != "" {
			result[trans.Key] = *trans.TargetText
		}
	}

	return e.writeJSON(result, destPath)
}

// ExportMerged exports translations merged with original content.
// Translated entries use TargetText, untranslated entries keep original.
func (e *Exporter) ExportMerged(originalContent []byte, translations []*models.TranslationWithSource, destPath string) error {
	// Parse original content
	var original map[string]interface{}
	if err := json.Unmarshal(originalContent, &original); err != nil {
		return fmt.Errorf("failed to parse original content: %w", err)
	}

	// Create translation lookup map
	translationMap := make(map[string]string)
	for _, trans := range translations {
		if trans.TargetText != nil && trans.Key != "" {
			translationMap[trans.Key] = *trans.TargetText
		}
	}

	// Build result with translations applied
	result := make(map[string]string)
	for key, value := range original {
		if translated, ok := translationMap[key]; ok {
			result[key] = translated
		} else if str, ok := value.(string); ok {
			result[key] = str
		}
	}

	return e.writeJSON(result, destPath)
}

// ExportTerms exports terms to a JSON file.
func (e *Exporter) ExportTerms(terms []*models.Term, destPath string) error {
	// Convert to exportable format
	result := make([]map[string]interface{}, 0, len(terms))

	for _, term := range terms {
		entry := map[string]interface{}{
			"scope":       term.Scope,
			"source_text": term.SourceText,
			"target_text": term.TargetText,
			"priority":    term.Priority,
		}

		if len(term.Tags) > 0 {
			entry["tags"] = term.Tags
		}
		if term.Context != nil {
			entry["context"] = *term.Context
		}
		if term.Notes != nil {
			entry["notes"] = *term.Notes
		}

		result = append(result, entry)
	}

	return e.writeJSONArray(result, destPath)
}

// writeJSON writes a map to a JSON file with sorted keys.
func (e *Exporter) writeJSON(data map[string]string, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build ordered map for JSON output
	orderedData := make(map[string]string, len(data))
	for _, key := range keys {
		orderedData[key] = data[key]
	}

	content, err := json.MarshalIndent(orderedData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// writeJSONArray writes an array to a JSON file.
func (e *Exporter) writeJSONArray(data []map[string]interface{}, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
