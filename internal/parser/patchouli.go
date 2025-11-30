package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// PatchouliParser parses Patchouli guidebook JSON files.
type PatchouliParser struct{}

// Compile-time check that PatchouliParser implements interfaces.Parser.
var _ interfaces.Parser = (*PatchouliParser)(nil)

// NewPatchouliParser creates a new Patchouli guidebook parser.
func NewPatchouliParser() *PatchouliParser {
	return &PatchouliParser{}
}

// patchouliPage represents a page in a Patchouli entry.
type patchouliPage struct {
	Type  string `json:"type"`
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// patchouliData represents the structure of a Patchouli JSON file.
type patchouliData struct {
	Name        string          `json:"name,omitempty"`
	Subtitle    string          `json:"subtitle,omitempty"`
	Description string          `json:"description,omitempty"`
	Pages       []patchouliPage `json:"pages,omitempty"`
}

// Parse extracts translation entries from Patchouli JSON content.
func (p *PatchouliParser) Parse(content []byte) ([]interfaces.ParsedEntry, error) {
	var data patchouliData
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	entries := []interfaces.ParsedEntry{}

	// Extract top-level translatable fields
	if data.Name != "" {
		entries = append(entries, interfaces.ParsedEntry{
			Key:  "name",
			Text: data.Name,
			Tags: []string{"patchouli"},
		})
	}

	if data.Subtitle != "" {
		entries = append(entries, interfaces.ParsedEntry{
			Key:  "subtitle",
			Text: data.Subtitle,
			Tags: []string{"patchouli"},
		})
	}

	if data.Description != "" {
		entries = append(entries, interfaces.ParsedEntry{
			Key:  "description",
			Text: data.Description,
			Tags: []string{"patchouli"},
		})
	}

	// Extract translatable fields from pages
	for i, page := range data.Pages {
		if page.Title != "" {
			entries = append(entries, interfaces.ParsedEntry{
				Key:  fmt.Sprintf("pages[%d].title", i),
				Text: page.Title,
				Tags: []string{"patchouli"},
			})
		}

		if page.Text != "" {
			entries = append(entries, interfaces.ParsedEntry{
				Key:  fmt.Sprintf("pages[%d].text", i),
				Text: page.Text,
				Tags: []string{"patchouli"},
			})
		}
	}

	return entries, nil
}

// Apply applies translations to the original Patchouli JSON content.
func (p *PatchouliParser) Apply(content []byte, translations map[string]string) ([]byte, error) {
	// Parse into generic map to preserve all fields
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Apply top-level translations
	if translation, exists := translations["name"]; exists {
		data["name"] = translation
	}

	if translation, exists := translations["subtitle"]; exists {
		data["subtitle"] = translation
	}

	if translation, exists := translations["description"]; exists {
		data["description"] = translation
	}

	// Apply page translations
	if pagesInterface, ok := data["pages"]; ok {
		pages, ok := pagesInterface.([]interface{})
		if ok {
			for i, pageInterface := range pages {
				page, ok := pageInterface.(map[string]interface{})
				if !ok {
					continue
				}

				// Apply title translation
				titleKey := fmt.Sprintf("pages[%d].title", i)
				if translation, exists := translations[titleKey]; exists {
					page["title"] = translation
				}

				// Apply text translation
				textKey := fmt.Sprintf("pages[%d].text", i)
				if translation, exists := translations[textKey]; exists {
					page["text"] = translation
				}

				pages[i] = page
			}
			data["pages"] = pages
		}
	}

	// Marshal with indentation for readability
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return result, nil
}

// SupportedTypes returns the pattern types this parser handles.
func (p *PatchouliParser) SupportedTypes() []string {
	return []string{"patchouli_category", "patchouli_entry"}
}
