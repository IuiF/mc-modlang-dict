package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

func TestNewExporter(t *testing.T) {
	exp := NewExporter()
	if exp == nil {
		t.Fatal("NewExporter() returned nil")
	}
}

func TestExporter_ExportJSON(t *testing.T) {
	exp := NewExporter()
	destDir := t.TempDir()

	translations := []*models.TranslationWithSource{
		{Key: "item.create.wrench", SourceText: "Wrench", Translation: models.Translation{TargetText: strPtr("レンチ")}},
		{Key: "block.create.gearbox", SourceText: "Gearbox", Translation: models.Translation{TargetText: strPtr("ギアボックス")}},
	}

	destPath := filepath.Join(destDir, "ja_jp.json")
	err := exp.ExportJSON(translations, destPath)
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	// Verify file contents
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	if result["item.create.wrench"] != "レンチ" {
		t.Errorf("ExportJSON() wrench = %v, want レンチ", result["item.create.wrench"])
	}
}

func TestExporter_ExportJSON_SkipUntranslated(t *testing.T) {
	exp := NewExporter()
	destDir := t.TempDir()

	translations := []*models.TranslationWithSource{
		{Key: "item.create.wrench", SourceText: "Wrench", Translation: models.Translation{TargetText: strPtr("レンチ")}},
		{Key: "block.create.gearbox", SourceText: "Gearbox", Translation: models.Translation{TargetText: nil}}, // Not translated
	}

	destPath := filepath.Join(destDir, "ja_jp.json")
	err := exp.ExportJSON(translations, destPath)
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("ExportJSON() got %d entries, want 1 (untranslated should be skipped)", len(result))
	}
}

func TestExporter_ExportMerged(t *testing.T) {
	exp := NewExporter()
	destDir := t.TempDir()

	originalContent := []byte(`{
  "item.create.wrench": "Wrench",
  "block.create.gearbox": "Gearbox",
  "tooltip.create.hint": "Hint"
}`)

	translations := []*models.TranslationWithSource{
		{Key: "item.create.wrench", SourceText: "Wrench", Translation: models.Translation{TargetText: strPtr("レンチ")}},
		{Key: "block.create.gearbox", SourceText: "Gearbox", Translation: models.Translation{TargetText: strPtr("ギアボックス")}},
	}

	destPath := filepath.Join(destDir, "ja_jp.json")
	err := exp.ExportMerged(originalContent, translations, destPath)
	if err != nil {
		t.Fatalf("ExportMerged() error = %v", err)
	}

	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	// Should have all 3 keys
	if len(result) != 3 {
		t.Errorf("ExportMerged() got %d entries, want 3", len(result))
	}

	// Translated keys should have Japanese text
	if result["item.create.wrench"] != "レンチ" {
		t.Errorf("ExportMerged() wrench = %v, want レンチ", result["item.create.wrench"])
	}

	// Untranslated key should keep original
	if result["tooltip.create.hint"] != "Hint" {
		t.Errorf("ExportMerged() hint = %v, want Hint", result["tooltip.create.hint"])
	}
}

func TestExporter_ExportTerms(t *testing.T) {
	exp := NewExporter()
	destDir := t.TempDir()

	terms := []*models.Term{
		{Scope: "global", SourceText: "Redstone", TargetText: "レッドストーン", Priority: 100},
		{Scope: "mod:create", SourceText: "Gearbox", TargetText: "ギアボックス", Priority: 200},
	}

	destPath := filepath.Join(destDir, "terms.json")
	err := exp.ExportTerms(terms, destPath)
	if err != nil {
		t.Fatalf("ExportTerms() error = %v", err)
	}

	// Verify file exists and is valid JSON
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Failed to parse output JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("ExportTerms() got %d terms, want 2", len(result))
	}
}

func strPtr(s string) *string {
	return &s
}
