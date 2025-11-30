package jar

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewExtractor(t *testing.T) {
	ext := NewExtractor()
	if ext == nil {
		t.Fatal("NewExtractor() returned nil")
	}
}

func TestExtractor_Extract(t *testing.T) {
	// Create a test JAR (ZIP) file
	jarPath := createTestJAR(t, map[string]string{
		"assets/testmod/lang/en_us.json": `{"item.testmod.test": "Test Item"}`,
		"fabric.mod.json":                `{"id": "testmod", "name": "Test Mod", "version": "1.0.0"}`,
	})
	defer os.Remove(jarPath)

	ext := NewExtractor()
	destDir := t.TempDir()

	result, err := ext.Extract(jarPath, destDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if result.ModID != "testmod" {
		t.Errorf("Extract() ModID = %v, want %v", result.ModID, "testmod")
	}

	// Check files were extracted
	langFile := filepath.Join(destDir, "assets", "testmod", "lang", "en_us.json")
	if _, err := os.Stat(langFile); os.IsNotExist(err) {
		t.Errorf("Extract() did not extract lang file")
	}
}

func TestExtractor_Extract_ForgeMod(t *testing.T) {
	// Create a Forge mod JAR
	modsToml := `
[[mods]]
modId = "create"
version = "0.5.1"
displayName = "Create"
authors = "simibubi"
`
	jarPath := createTestJAR(t, map[string]string{
		"META-INF/mods.toml":             modsToml,
		"assets/create/lang/en_us.json":  `{"item.create.wrench": "Wrench"}`,
	})
	defer os.Remove(jarPath)

	ext := NewExtractor()
	destDir := t.TempDir()

	result, err := ext.Extract(jarPath, destDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if result.ModID != "create" {
		t.Errorf("Extract() ModID = %v, want %v", result.ModID, "create")
	}

	if result.DisplayName != "Create" {
		t.Errorf("Extract() DisplayName = %v, want %v", result.DisplayName, "Create")
	}
}

func TestExtractor_DetectModID_Fabric(t *testing.T) {
	ext := NewExtractor()

	content := []byte(`{
		"id": "botania",
		"name": "Botania",
		"version": "1.2.3",
		"authors": ["Vazkii"]
	}`)

	info, err := ext.detectFabricMod(content)
	if err != nil {
		t.Fatalf("detectFabricMod() error = %v", err)
	}

	if info.ModID != "botania" {
		t.Errorf("ModID = %v, want %v", info.ModID, "botania")
	}
	if info.DisplayName != "Botania" {
		t.Errorf("DisplayName = %v, want %v", info.DisplayName, "Botania")
	}
}

func TestExtractor_DetectModID_Forge(t *testing.T) {
	ext := NewExtractor()

	content := []byte(`
[[mods]]
modId = "thermal"
version = "1.0.0"
displayName = "Thermal Expansion"
authors = "TeamCoFH"
description = "Thermal tech mod"
`)

	info, err := ext.detectForgeMod(content)
	if err != nil {
		t.Fatalf("detectForgeMod() error = %v", err)
	}

	if info.ModID != "thermal" {
		t.Errorf("ModID = %v, want %v", info.ModID, "thermal")
	}
	if info.DisplayName != "Thermal Expansion" {
		t.Errorf("DisplayName = %v, want %v", info.DisplayName, "Thermal Expansion")
	}
}

func TestExtractor_ListLangFiles(t *testing.T) {
	jarPath := createTestJAR(t, map[string]string{
		"assets/testmod/lang/en_us.json": `{}`,
		"assets/testmod/lang/ja_jp.json": `{}`,
		"assets/testmod/lang/zh_cn.json": `{}`,
		"assets/testmod/textures/item.png": "binary",
		"fabric.mod.json": `{"id": "testmod"}`,
	})
	defer os.Remove(jarPath)

	ext := NewExtractor()
	destDir := t.TempDir()

	result, err := ext.Extract(jarPath, destDir)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	if len(result.LangFiles) != 3 {
		t.Errorf("Extract() found %d lang files, want 3", len(result.LangFiles))
	}
}

// createTestJAR creates a temporary JAR file with the given contents.
func createTestJAR(t *testing.T, files map[string]string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-*.jar")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for name, content := range files {
		w, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("Failed to create zip entry: %v", err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write zip entry: %v", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("Failed to close zip: %v", err)
	}

	if err := os.WriteFile(tmpFile.Name(), buf.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to write JAR file: %v", err)
	}

	return tmpFile.Name()
}
