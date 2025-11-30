package parser

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// JSONLangParser parses Minecraft JSON language files.
type JSONLangParser struct{}

// Compile-time check that JSONLangParser implements interfaces.Parser.
var _ interfaces.Parser = (*JSONLangParser)(nil)

// NewJSONLangParser creates a new JSON language file parser.
func NewJSONLangParser() *JSONLangParser {
	return &JSONLangParser{}
}

// Parse extracts translation entries from JSON language file content.
func (p *JSONLangParser) Parse(content []byte) ([]interfaces.ParsedEntry, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	entries := make([]interfaces.ParsedEntry, 0, len(data))

	// Get sorted keys for consistent ordering
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := data[key]

		// Only process string values
		text, ok := value.(string)
		if !ok {
			continue
		}

		entry := interfaces.ParsedEntry{
			Key:  key,
			Text: text,
			Tags: detectTags(key),
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// Apply applies translations to the original JSON content.
func (p *JSONLangParser) Apply(content []byte, translations map[string]string) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Apply translations
	for key, translation := range translations {
		if _, exists := data[key]; exists {
			data[key] = translation
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
func (p *JSONLangParser) SupportedTypes() []string {
	return []string{"lang"}
}

// detectTags extracts tags from a translation key.
// Keys typically follow the format: type.modid.name
// Examples: item.create.wrench, block.create.gearbox
func detectTags(key string) []string {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) < 1 {
		return nil
	}

	prefix := strings.ToLower(parts[0])

	// Known Minecraft translation key prefixes
	knownTypes := map[string]bool{
		"item":          true,
		"block":         true,
		"entity":        true,
		"tooltip":       true,
		"enchantment":   true,
		"effect":        true,
		"potion":        true,
		"biome":         true,
		"advancement":   true,
		"container":     true,
		"commands":      true,
		"death":         true,
		"key":           true,
		"gui":           true,
		"chat":          true,
		"subtitles":     true,
		"jei":           true,
		"config":        true,
		"message":       true,
		"itemGroup":     true,
		"creativetab":   true,
		"fluid":         true,
		"sound":         true,
		"stat":          true,
		"attribute":     true,
		"recipe":        true,
		"painting":      true,
		"record":        true,
		"trim_material": true,
		"trim_pattern":  true,
	}

	if knownTypes[prefix] {
		return []string{prefix}
	}

	return nil
}
