package parser

import (
	"encoding/json"
	"fmt"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// MantleBookParser parses Mantle Book JSON files (used by Tinkers' Construct).
type MantleBookParser struct{}

// Compile-time check that MantleBookParser implements interfaces.Parser.
var _ interfaces.Parser = (*MantleBookParser)(nil)

// NewMantleBookParser creates a new Mantle Book parser.
func NewMantleBookParser() *MantleBookParser {
	return &MantleBookParser{}
}

// Parse extracts translation entries from Mantle Book JSON content.
func (p *MantleBookParser) Parse(content []byte) ([]interfaces.ParsedEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	entries := []interfaces.ParsedEntry{}

	// Extract title field
	if title, ok := data["title"].(string); ok && title != "" {
		entries = append(entries, interfaces.ParsedEntry{
			Key:  "title",
			Text: title,
			Tags: []string{"mantle_book"},
		})
	}

	// Extract text field - can be string or array
	if textValue, ok := data["text"]; ok {
		switch v := textValue.(type) {
		case string:
			// Text as string
			if v != "" {
				entries = append(entries, interfaces.ParsedEntry{
					Key:  "text",
					Text: v,
					Tags: []string{"mantle_book"},
				})
			}
		case []interface{}:
			// Text as array of objects
			for i, item := range v {
				if textObj, ok := item.(map[string]interface{}); ok {
					if text, ok := textObj["text"].(string); ok && text != "" {
						entries = append(entries, interfaces.ParsedEntry{
							Key:  fmt.Sprintf("text[%d].text", i),
							Text: text,
							Tags: []string{"mantle_book"},
						})
					}
				}
			}
		}
	}

	// Extract string arrays: properties, entity, block
	entries = append(entries, extract_string_array(data, "properties")...)
	entries = append(entries, extract_string_array(data, "entity")...)
	entries = append(entries, extract_string_array(data, "block")...)

	return entries, nil
}

// extract_string_array extracts entries from a string array field.
func extract_string_array(data map[string]interface{}, fieldName string) []interfaces.ParsedEntry {
	entries := []interfaces.ParsedEntry{}

	if arr, ok := data[fieldName].([]interface{}); ok {
		for i, item := range arr {
			if str, ok := item.(string); ok && str != "" {
				entries = append(entries, interfaces.ParsedEntry{
					Key:  fmt.Sprintf("%s[%d]", fieldName, i),
					Text: str,
					Tags: []string{"mantle_book"},
				})
			}
		}
	}

	return entries
}

// Apply applies translations to the original Mantle Book JSON content.
func (p *MantleBookParser) Apply(content []byte, translations map[string]string) ([]byte, error) {
	// Parse into generic map to preserve all fields
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Apply title translation
	if translation, exists := translations["title"]; exists {
		data["title"] = translation
	}

	// Apply text field translations
	if textValue, ok := data["text"]; ok {
		switch v := textValue.(type) {
		case string:
			// Text as string
			if translation, exists := translations["text"]; exists {
				data["text"] = translation
			}
		case []interface{}:
			// Text as array of objects
			for i, item := range v {
				if textObj, ok := item.(map[string]interface{}); ok {
					key := fmt.Sprintf("text[%d].text", i)
					if translation, exists := translations[key]; exists {
						textObj["text"] = translation
					}
					v[i] = textObj
				}
			}
			data["text"] = v
		}
	}

	// Apply string array translations: properties, entity, block
	apply_string_array_translations(data, "properties", translations)
	apply_string_array_translations(data, "entity", translations)
	apply_string_array_translations(data, "block", translations)

	// Marshal with indentation for readability
	result, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return result, nil
}

// apply_string_array_translations applies translations to a string array field.
func apply_string_array_translations(data map[string]interface{}, fieldName string, translations map[string]string) {
	if arr, ok := data[fieldName].([]interface{}); ok {
		for i := range arr {
			key := fmt.Sprintf("%s[%d]", fieldName, i)
			if translation, exists := translations[key]; exists {
				arr[i] = translation
			}
		}
		data[fieldName] = arr
	}
}

// SupportedTypes returns the pattern types this parser handles.
func (p *MantleBookParser) SupportedTypes() []string {
	return []string{"mantle_book"}
}
