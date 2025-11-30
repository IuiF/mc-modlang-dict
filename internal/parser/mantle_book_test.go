package parser

import (
	"testing"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// TestMantleBookParser_Parse tests the Parse method with various Mantle Book JSON structures.
func TestMantleBookParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantCount  int
		wantErrors bool
		checkEntry func(t *testing.T, entries []interfaces.ParsedEntry)
	}{
		{
			name: "basic page with title and text array",
			content: `{
				"title": "Focus",
				"text": [
					{
						"text": "This book is focused on tool building."
					},
					{
						"text": "If you are unsure of how to use the technology.",
						"paragraph": true
					}
				]
			}`,
			wantCount: 3,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				// Check title
				if entries[0].Key != "title" {
					t.Errorf("entries[0].Key = %q, want %q", entries[0].Key, "title")
				}
				if entries[0].Text != "Focus" {
					t.Errorf("entries[0].Text = %q, want %q", entries[0].Text, "Focus")
				}
				// Check first text
				if entries[1].Key != "text[0].text" {
					t.Errorf("entries[1].Key = %q, want %q", entries[1].Key, "text[0].text")
				}
				// Check second text
				if entries[2].Key != "text[1].text" {
					t.Errorf("entries[2].Key = %q, want %q", entries[2].Key, "text[1].text")
				}
			},
		},
		{
			name: "tool page with properties",
			content: `{
				"tool": "tconstruct:pickaxe",
				"text": [
					{
						"text": "The Pickaxe is a precise mining tool."
					}
				],
				"properties": [
					"+0.5 Attack Damage",
					"1.2 Attack Speed",
					"Piercing I"
				]
			}`,
			wantCount: 4,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				// Check text entry
				if entries[0].Key != "text[0].text" {
					t.Errorf("entries[0].Key = %q, want %q", entries[0].Key, "text[0].text")
				}
				// Check properties
				if entries[1].Key != "properties[0]" {
					t.Errorf("entries[1].Key = %q, want %q", entries[1].Key, "properties[0]")
				}
				if entries[1].Text != "+0.5 Attack Damage" {
					t.Errorf("entries[1].Text = %q, want %q", entries[1].Text, "+0.5 Attack Damage")
				}
			},
		},
		{
			name: "fluid page with text string and entity/block arrays",
			content: `{
				"title": "Water",
				"text": "Water is not the most useful fluid.",
				"entity": [
					"2 Water Damage",
					"Extinguishes Fire"
				],
				"block": [
					"Breaks blocks broken by water"
				]
			}`,
			wantCount: 5,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				// Check title
				if entries[0].Key != "title" || entries[0].Text != "Water" {
					t.Errorf("title entry incorrect: got (%q, %q)", entries[0].Key, entries[0].Text)
				}
				// Check text as string
				if entries[1].Key != "text" || entries[1].Text != "Water is not the most useful fluid." {
					t.Errorf("text entry incorrect: got (%q, %q)", entries[1].Key, entries[1].Text)
				}
				// Check entity array
				if entries[2].Key != "entity[0]" || entries[2].Text != "2 Water Damage" {
					t.Errorf("entity[0] entry incorrect: got (%q, %q)", entries[2].Key, entries[2].Text)
				}
				// Check block array
				if entries[4].Key != "block[0]" || entries[4].Text != "Breaks blocks broken by water" {
					t.Errorf("block[0] entry incorrect: got (%q, %q)", entries[4].Key, entries[4].Text)
				}
			},
		},
		{
			name: "empty JSON object",
			content: `{}`,
			wantCount: 0,
		},
		{
			name:       "invalid JSON",
			content:    `{"invalid": json}`,
			wantErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMantleBookParser()
			entries, err := parser.Parse([]byte(tt.content))

			if tt.wantErrors {
				if err == nil {
					t.Errorf("Parse() error = nil, wantErrors %v", tt.wantErrors)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error = %v", err)
			}

			if len(entries) != tt.wantCount {
				t.Errorf("Parse() returned %d entries, want %d", len(entries), tt.wantCount)
			}

			if tt.checkEntry != nil && len(entries) > 0 {
				tt.checkEntry(t, entries)
			}

			// All entries should have mantle_book tag
			for i, entry := range entries {
				found := false
				for _, tag := range entry.Tags {
					if tag == "mantle_book" {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("entries[%d] missing 'mantle_book' tag", i)
				}
			}
		})
	}
}

// TestMantleBookParser_Apply tests the Apply method.
func TestMantleBookParser_Apply(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		translations map[string]string
		wantContains []string
	}{
		{
			name: "apply translations to basic page",
			content: `{
				"title": "Focus",
				"text": [
					{"text": "This book is focused on tool building."}
				]
			}`,
			translations: map[string]string{
				"title":         "フォーカス",
				"text[0].text":  "この本はツール作成に焦点を当てています。",
			},
			wantContains: []string{
				`"title": "フォーカス"`,
				`"この本はツール作成に焦点を当てています。"`,
			},
		},
		{
			name: "apply translations to tool page",
			content: `{
				"tool": "tconstruct:pickaxe",
				"properties": ["+0.5 Attack Damage"]
			}`,
			translations: map[string]string{
				"properties[0]": "+0.5 攻撃ダメージ",
			},
			wantContains: []string{
				`"tool": "tconstruct:pickaxe"`,
				`"+0.5 攻撃ダメージ"`,
			},
		},
		{
			name: "apply translations to fluid page with string text",
			content: `{
				"title": "Water",
				"text": "Water is not the most useful fluid."
			}`,
			translations: map[string]string{
				"title": "水",
				"text":  "水はあまり便利な液体ではありません。",
			},
			wantContains: []string{
				`"title": "水"`,
				`"text": "水はあまり便利な液体ではありません。"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMantleBookParser()
			result, err := parser.Apply([]byte(tt.content), tt.translations)

			if err != nil {
				t.Fatalf("Apply() unexpected error = %v", err)
			}

			resultStr := string(result)
			for _, want := range tt.wantContains {
				if !contains(resultStr, want) {
					t.Errorf("Apply() result does not contain %q\nGot: %s", want, resultStr)
				}
			}
		})
	}
}

// TestMantleBookParser_SupportedTypes tests the SupportedTypes method.
func TestMantleBookParser_SupportedTypes(t *testing.T) {
	parser := NewMantleBookParser()
	types := parser.SupportedTypes()

	expectedTypes := []string{"mantle_book"}
	if len(types) != len(expectedTypes) {
		t.Errorf("SupportedTypes() returned %d types, want %d", len(types), len(expectedTypes))
	}

	for i, expected := range expectedTypes {
		if i >= len(types) || types[i] != expected {
			t.Errorf("SupportedTypes()[%d] = %q, want %q", i, types[i], expected)
		}
	}
}

// TestMantleBookParser_ImplementsInterface verifies interface implementation.
func TestMantleBookParser_ImplementsInterface(t *testing.T) {
	var _ interfaces.Parser = (*MantleBookParser)(nil)
}

// TestMantleBookParser_RealWorldData tests with actual Tinkers' Construct JSON data.
func TestMantleBookParser_RealWorldData(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantCount  int
		checkEntry func(t *testing.T, entries []interfaces.ParsedEntry)
	}{
		{
			name: "actual foreward.json from Tinkers' Construct",
			content: `{
				"title": "Foreword",
				"text": [
					{
						"text": "During my research into tool making, I noticed several authors have created partial lists of materials but there was no one complete list. This means in order to properly compare all the materials, it is necessary to switch between several volumes. This book is intended to solve that problem by indexing known materials from several different volumes into a single book."
					},
					{
						"text": "I could not have created such a grand book alone. This book also includes lists of materials collected by both Thruul M'Gon and Nemea. Just don't tell either of them that I spoke with both of them.",
						"paragraph": true
					}
				]
			}`,
			wantCount: 3,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				if entries[0].Text != "Foreword" {
					t.Errorf("Expected title 'Foreword', got %q", entries[0].Text)
				}
			},
		},
		{
			name: "actual tconstruct_pickaxe.json from Tinkers' Construct",
			content: `{
				"tool": "tconstruct:pickaxe",
				"text": [
					{
						"text": "The Pickaxe is a precise mining tool, effective on stone, metal, and ores."
					},
					{
						"text": "Expanders alternate between increasing depth and increasing height.",
						"paragraph": true
					}
				],
				"properties": [
					"+0.5 Attack Damage",
					"1.2 Attack Speed",
					"Piercing I"
				]
			}`,
			wantCount: 5,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				// Verify we have text entries and properties
				hasText := false
				hasProperties := false
				for _, entry := range entries {
					if entry.Key == "text[0].text" {
						hasText = true
					}
					if entry.Key == "properties[0]" {
						hasProperties = true
					}
				}
				if !hasText {
					t.Error("Expected to find text entries")
				}
				if !hasProperties {
					t.Error("Expected to find properties entries")
				}
			},
		},
		{
			name: "actual tconstruct_water.json from Tinkers' Construct",
			content: `{
				"title": "Water",
				"text": "Water is not the most useful fluid, but can protect you from fire in hot environments.",
				"entity": [
					"2 Water Damage",
					"Extinguishes Fire"
				],
				"block": [
					"Breaks blocks broken by water"
				]
			}`,
			wantCount: 5,
			checkEntry: func(t *testing.T, entries []interfaces.ParsedEntry) {
				// Verify all field types are present
				hasTitle := false
				hasText := false
				hasEntity := false
				hasBlock := false
				for _, entry := range entries {
					switch entry.Key {
					case "title":
						hasTitle = true
					case "text":
						hasText = true
					case "entity[0]":
						hasEntity = true
					case "block[0]":
						hasBlock = true
					}
				}
				if !hasTitle || !hasText || !hasEntity || !hasBlock {
					t.Errorf("Missing expected fields: title=%v, text=%v, entity=%v, block=%v",
						hasTitle, hasText, hasEntity, hasBlock)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMantleBookParser()
			entries, err := parser.Parse([]byte(tt.content))

			if err != nil {
				t.Fatalf("Parse() unexpected error = %v", err)
			}

			if len(entries) != tt.wantCount {
				t.Errorf("Parse() returned %d entries, want %d", len(entries), tt.wantCount)
			}

			if tt.checkEntry != nil {
				tt.checkEntry(t, entries)
			}
		})
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
