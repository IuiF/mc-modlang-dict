package parser

import (
	"testing"
)

func TestJSONLangParser_Parse(t *testing.T) {
	parser := NewJSONLangParser()

	tests := []struct {
		name      string
		content   string
		wantCount int
		wantKeys  []string
	}{
		{
			name: "basic",
			content: `{
				"item.create.wrench": "Wrench",
				"block.create.gearbox": "Gearbox"
			}`,
			wantCount: 2,
			wantKeys:  []string{"item.create.wrench", "block.create.gearbox"},
		},
		{
			name:      "empty",
			content:   `{}`,
			wantCount: 0,
			wantKeys:  []string{},
		},
		{
			name: "with_formatting",
			content: `{
				"tooltip.create.speed": "Speed: %s RPM",
				"message.create.welcome": "Welcome, %1$s!"
			}`,
			wantCount: 2,
			wantKeys:  []string{"tooltip.create.speed", "message.create.welcome"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Errorf("Parse() error = %v", err)
				return
			}

			if len(entries) != tt.wantCount {
				t.Errorf("Parse() got %d entries, want %d", len(entries), tt.wantCount)
			}

			keyMap := make(map[string]bool)
			for _, entry := range entries {
				keyMap[entry.Key] = true
			}

			for _, wantKey := range tt.wantKeys {
				if !keyMap[wantKey] {
					t.Errorf("Parse() missing key %q", wantKey)
				}
			}
		})
	}
}

func TestJSONLangParser_Parse_Invalid(t *testing.T) {
	parser := NewJSONLangParser()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "invalid_json",
			content: `{invalid}`,
		},
		{
			name:    "array_instead_of_object",
			content: `["item1", "item2"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse([]byte(tt.content))
			if err == nil {
				t.Error("Parse() expected error for invalid content")
			}
		})
	}
}

func TestJSONLangParser_Apply(t *testing.T) {
	parser := NewJSONLangParser()

	original := `{
  "item.create.wrench": "Wrench",
  "block.create.gearbox": "Gearbox"
}`

	translations := map[string]string{
		"item.create.wrench":   "レンチ",
		"block.create.gearbox": "ギアボックス",
	}

	result, err := parser.Apply([]byte(original), translations)
	if err != nil {
		t.Errorf("Apply() error = %v", err)
		return
	}

	// Parse result to verify
	entries, err := parser.Parse(result)
	if err != nil {
		t.Errorf("Parse result error = %v", err)
		return
	}

	for _, entry := range entries {
		expected, ok := translations[entry.Key]
		if ok && entry.Text != expected {
			t.Errorf("Apply() key %q = %q, want %q", entry.Key, entry.Text, expected)
		}
	}
}

func TestJSONLangParser_SupportedTypes(t *testing.T) {
	parser := NewJSONLangParser()

	types := parser.SupportedTypes()
	if len(types) == 0 {
		t.Error("SupportedTypes() returned empty slice")
	}

	found := false
	for _, typ := range types {
		if typ == "lang" {
			found = true
			break
		}
	}

	if !found {
		t.Error("SupportedTypes() does not include 'lang'")
	}
}

func TestJSONLangParser_TagDetection(t *testing.T) {
	parser := NewJSONLangParser()

	content := `{
		"item.create.wrench": "Wrench",
		"block.create.gearbox": "Gearbox",
		"entity.create.contraption": "Contraption",
		"tooltip.create.hint": "Hint"
	}`

	entries, err := parser.Parse([]byte(content))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	expectedTags := map[string]string{
		"item.create.wrench":        "item",
		"block.create.gearbox":      "block",
		"entity.create.contraption": "entity",
		"tooltip.create.hint":       "tooltip",
	}

	for _, entry := range entries {
		expectedTag := expectedTags[entry.Key]
		found := false
		for _, tag := range entry.Tags {
			if tag == expectedTag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Entry %q missing expected tag %q, got %v", entry.Key, expectedTag, entry.Tags)
		}
	}
}
