// Package jar provides functionality for extracting and analyzing Minecraft mod JAR files.
package jar

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ExtractResult contains information extracted from a mod JAR file.
type ExtractResult struct {
	ModID       string            // Detected mod ID
	DisplayName string            // Display name of the mod
	Version     string            // Mod version
	MCVersion   string            // Minecraft version (if available)
	Loader      string            // Mod loader (forge, fabric, neoforge, quilt)
	Authors     []string          // Mod authors
	Description string            // Mod description
	LangFiles   []string          // Paths to extracted lang files
	ExtractDir  string            // Directory where files were extracted
	Metadata    map[string]string // Additional metadata
}

// Extractor handles JAR file extraction and mod detection.
type Extractor struct{}

// NewExtractor creates a new JAR extractor.
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts a JAR file and detects mod information.
func (e *Extractor) Extract(jarPath, destDir string) (*ExtractResult, error) {
	reader, err := zip.OpenReader(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open JAR: %w", err)
	}
	defer reader.Close()

	result := &ExtractResult{
		ExtractDir: destDir,
		LangFiles:  make([]string, 0),
		Metadata:   make(map[string]string),
	}

	// First pass: detect mod info and collect file list
	var fabricModJSON, modsToml []byte
	var isNeoForge bool

	for _, file := range reader.File {
		switch file.Name {
		case "fabric.mod.json":
			fabricModJSON, err = readZipFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read fabric.mod.json: %w", err)
			}
		case "META-INF/neoforge.mods.toml":
			// NeoForge uses neoforge.mods.toml (prioritize over mods.toml)
			modsToml, err = readZipFile(file)
			if err != nil {
				return nil, fmt.Errorf("failed to read neoforge.mods.toml: %w", err)
			}
			isNeoForge = true
		case "META-INF/mods.toml":
			// Only use mods.toml if we haven't found neoforge.mods.toml
			if modsToml == nil {
				modsToml, err = readZipFile(file)
				if err != nil {
					return nil, fmt.Errorf("failed to read mods.toml: %w", err)
				}
			}
		case "quilt.mod.json":
			// Quilt uses similar format to Fabric
			if fabricModJSON == nil {
				fabricModJSON, err = readZipFile(file)
				if err != nil {
					return nil, fmt.Errorf("failed to read quilt.mod.json: %w", err)
				}
				result.Loader = "quilt"
			}
		}
	}

	// Detect mod info
	if fabricModJSON != nil {
		info, err := e.detectFabricMod(fabricModJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to parse fabric.mod.json: %w", err)
		}
		result.ModID = info.ModID
		result.DisplayName = info.DisplayName
		result.Version = info.Version
		result.Authors = info.Authors
		result.Description = info.Description
		if result.Loader == "" {
			result.Loader = "fabric"
		}
	} else if modsToml != nil {
		info, err := e.detectForgeMod(modsToml)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mods.toml: %w", err)
		}
		result.ModID = info.ModID
		result.DisplayName = info.DisplayName
		result.Version = info.Version
		result.Authors = info.Authors
		result.Description = info.Description
		if isNeoForge {
			result.Loader = "neoforge"
		} else {
			result.Loader = "forge"
		}
	}

	// Second pass: extract files
	for _, file := range reader.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		destPath := filepath.Join(destDir, file.Name)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}

		// Extract file
		if err := extractFile(file, destPath); err != nil {
			return nil, fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}

		// Track lang files
		if isLangFile(file.Name) {
			result.LangFiles = append(result.LangFiles, destPath)
		}
	}

	return result, nil
}

// fabricModInfo represents fabric.mod.json structure.
type fabricModInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Authors     []interface{} `json:"authors"` // Can be string or object
	Description string   `json:"description"`
}

// detectFabricMod parses fabric.mod.json content.
func (e *Extractor) detectFabricMod(content []byte) (*ExtractResult, error) {
	var info fabricModInfo
	if err := json.Unmarshal(content, &info); err != nil {
		return nil, err
	}

	authors := make([]string, 0, len(info.Authors))
	for _, author := range info.Authors {
		switch v := author.(type) {
		case string:
			authors = append(authors, v)
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				authors = append(authors, name)
			}
		}
	}

	return &ExtractResult{
		ModID:       info.ID,
		DisplayName: info.Name,
		Version:     info.Version,
		Authors:     authors,
		Description: info.Description,
	}, nil
}

// forgeModInfo represents mods.toml structure.
type forgeModInfo struct {
	Mods []struct {
		ModID       string `toml:"modId"`
		Version     string `toml:"version"`
		DisplayName string `toml:"displayName"`
		Authors     string `toml:"authors"`
		Description string `toml:"description"`
	} `toml:"mods"`
}

// detectForgeMod parses mods.toml content.
func (e *Extractor) detectForgeMod(content []byte) (*ExtractResult, error) {
	var info forgeModInfo
	if _, err := toml.Decode(string(content), &info); err != nil {
		return nil, err
	}

	if len(info.Mods) == 0 {
		return nil, fmt.Errorf("no mods found in mods.toml")
	}

	mod := info.Mods[0]
	authors := []string{}
	if mod.Authors != "" {
		authors = []string{mod.Authors}
	}

	return &ExtractResult{
		ModID:       mod.ModID,
		DisplayName: mod.DisplayName,
		Version:     mod.Version,
		Authors:     authors,
		Description: mod.Description,
	}, nil
}

// readZipFile reads the content of a zip file entry.
func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

// extractFile extracts a single file from the zip archive.
func extractFile(file *zip.File, destPath string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// isLangFile checks if a file path is a language file.
func isLangFile(path string) bool {
	// Standard Minecraft lang file pattern: assets/{mod_id}/lang/{lang}.json
	if !strings.HasSuffix(path, ".json") {
		return false
	}

	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return false
	}

	return parts[0] == "assets" && parts[2] == "lang"
}
