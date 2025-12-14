package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// Compile-time check that LegacyLangParser implements interfaces.Parser.
var _ interfaces.Parser = (*LegacyLangParser)(nil)

// LegacyLangParser handles .lang files (Minecraft 1.12.2 and earlier format)
type LegacyLangParser struct{}

// NewLegacyLangParser creates a new parser for .lang files
func NewLegacyLangParser() *LegacyLangParser {
	return &LegacyLangParser{}
}

// Parse parses .lang file content
func (p *LegacyLangParser) Parse(content []byte) ([]interfaces.ParsedEntry, error) {
	var entries []interfaces.ParsedEntry
	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			continue
		}

		entries = append(entries, interfaces.ParsedEntry{
			Key:  key,
			Text: value,
			Tags: detectTags(key),
		})
	}

	return entries, nil
}

// Apply applies translations to the original .lang content.
func (p *LegacyLangParser) Apply(content []byte, translations map[string]string) ([]byte, error) {
	var result bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Keep empty lines and comments as-is
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		// Parse key=value format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}

		key := strings.TrimSpace(parts[0])

		// Apply translation if available
		if trans, ok := translations[key]; ok {
			result.WriteString(fmt.Sprintf("%s=%s\n", key, trans))
		} else {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.Bytes(), nil
}

// SupportedTypes returns the pattern types this parser handles.
func (p *LegacyLangParser) SupportedTypes() []string {
	return []string{"lang"}
}
