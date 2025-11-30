package parser

import (
	"testing"
)

func TestPatchouliParser_Parse_Category(t *testing.T) {
	parser := NewPatchouliParser()

	tests := []struct {
		name      string
		content   string
		wantCount int
		wantKeys  []string
		wantTexts map[string]string
	}{
		{
			name: "basic_category",
			content: `{
				"name": "Getting Started",
				"description": "Learn the basics",
				"icon": "minecraft:book"
			}`,
			wantCount: 2,
			wantKeys:  []string{"name", "description"},
			wantTexts: map[string]string{
				"name":        "Getting Started",
				"description": "Learn the basics",
			},
		},
		{
			name: "category_with_formatting_codes",
			content: `{
				"name": "$(item)Advanced Topics/$",
				"description": "Deep dive into $(l)advanced features/$ and techniques"
			}`,
			wantCount: 2,
			wantKeys:  []string{"name", "description"},
			wantTexts: map[string]string{
				"name":        "$(item)Advanced Topics/$",
				"description": "Deep dive into $(l)advanced features/$ and techniques",
			},
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

			keyMap := make(map[string]string)
			for _, entry := range entries {
				keyMap[entry.Key] = entry.Text
			}

			for _, wantKey := range tt.wantKeys {
				if _, exists := keyMap[wantKey]; !exists {
					t.Errorf("Parse() missing key %q", wantKey)
				}
			}

			for key, wantText := range tt.wantTexts {
				if gotText, exists := keyMap[key]; exists {
					if gotText != wantText {
						t.Errorf("Parse() key %q = %q, want %q", key, gotText, wantText)
					}
				}
			}
		})
	}
}

func TestPatchouliParser_Parse_Entry(t *testing.T) {
	parser := NewPatchouliParser()

	tests := []struct {
		name      string
		content   string
		wantCount int
		wantKeys  []string
		wantTexts map[string]string
	}{
		{
			name: "basic_entry",
			content: `{
				"name": "Welcome",
				"icon": "minecraft:book",
				"category": "patchouli:getting_started",
				"pages": [
					{
						"type": "text",
						"text": "Welcome to the guide!"
					}
				]
			}`,
			wantCount: 2,
			wantKeys:  []string{"name", "pages[0].text"},
			wantTexts: map[string]string{
				"name":          "Welcome",
				"pages[0].text": "Welcome to the guide!",
			},
		},
		{
			name: "entry_with_multiple_pages",
			content: `{
				"name": "Advanced Guide",
				"icon": "minecraft:diamond",
				"category": "patchouli:advanced",
				"pages": [
					{
						"type": "text",
						"title": "Introduction",
						"text": "This is an introduction."
					},
					{
						"type": "text",
						"title": "Details",
						"text": "Here are the details."
					},
					{
						"type": "crafting",
						"recipe": "minecraft:diamond_pickaxe",
						"text": "Craft this item."
					}
				]
			}`,
			wantCount: 6,
			wantKeys: []string{
				"name",
				"pages[0].title",
				"pages[0].text",
				"pages[1].title",
				"pages[1].text",
				"pages[2].text",
			},
			wantTexts: map[string]string{
				"name":           "Advanced Guide",
				"pages[0].title": "Introduction",
				"pages[0].text":  "This is an introduction.",
				"pages[1].title": "Details",
				"pages[1].text":  "Here are the details.",
				"pages[2].text":  "Craft this item.",
			},
		},
		{
			name: "entry_with_subtitle",
			content: `{
				"name": "Special Entry",
				"subtitle": "A subtitle here",
				"icon": "minecraft:emerald",
				"category": "patchouli:special",
				"pages": [
					{
						"type": "text",
						"text": "Content here."
					}
				]
			}`,
			wantCount: 3,
			wantKeys:  []string{"name", "subtitle", "pages[0].text"},
			wantTexts: map[string]string{
				"name":          "Special Entry",
				"subtitle":      "A subtitle here",
				"pages[0].text": "Content here.",
			},
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

			keyMap := make(map[string]string)
			for _, entry := range entries {
				keyMap[entry.Key] = entry.Text
			}

			for _, wantKey := range tt.wantKeys {
				if _, exists := keyMap[wantKey]; !exists {
					t.Errorf("Parse() missing key %q", wantKey)
				}
			}

			for key, wantText := range tt.wantTexts {
				if gotText, exists := keyMap[key]; exists {
					if gotText != wantText {
						t.Errorf("Parse() key %q = %q, want %q", key, gotText, wantText)
					}
				}
			}
		})
	}
}

func TestPatchouliParser_Parse_Invalid(t *testing.T) {
	parser := NewPatchouliParser()

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
		{
			name:    "empty_json",
			content: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := parser.Parse([]byte(tt.content))
			// empty_jsonの場合はエラーではなく、空のエントリを返す
			if tt.name == "empty_json" {
				if err != nil {
					t.Errorf("Parse() unexpected error for empty JSON: %v", err)
				}
				if len(entries) != 0 {
					t.Errorf("Parse() expected 0 entries for empty JSON, got %d", len(entries))
				}
			} else {
				if err == nil {
					t.Error("Parse() expected error for invalid content")
				}
			}
		})
	}
}

func TestPatchouliParser_Apply_Category(t *testing.T) {
	parser := NewPatchouliParser()

	original := `{
  "name": "Getting Started",
  "description": "Learn the basics",
  "icon": "minecraft:book"
}`

	translations := map[string]string{
		"name":        "はじめに",
		"description": "基本を学ぶ",
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

	entryMap := make(map[string]string)
	for _, entry := range entries {
		entryMap[entry.Key] = entry.Text
	}

	for key, expected := range translations {
		if got, exists := entryMap[key]; !exists {
			t.Errorf("Apply() missing key %q", key)
		} else if got != expected {
			t.Errorf("Apply() key %q = %q, want %q", key, got, expected)
		}
	}

	// Verify icon is preserved
	if got := entryMap["icon"]; got != "" {
		// icon should not be in parsed entries as it's not translatable
		t.Errorf("Apply() icon should not be in parsed entries")
	}
}

func TestPatchouliParser_Apply_Entry(t *testing.T) {
	parser := NewPatchouliParser()

	original := `{
  "name": "Welcome",
  "icon": "minecraft:book",
  "category": "patchouli:getting_started",
  "pages": [
    {
      "type": "text",
      "text": "Welcome to the guide!"
    },
    {
      "type": "text",
      "title": "More Info",
      "text": "Additional information here."
    }
  ]
}`

	translations := map[string]string{
		"name":           "ようこそ",
		"pages[0].text":  "ガイドへようこそ！",
		"pages[1].title": "追加情報",
		"pages[1].text":  "ここに追加情報があります。",
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

	entryMap := make(map[string]string)
	for _, entry := range entries {
		entryMap[entry.Key] = entry.Text
	}

	for key, expected := range translations {
		if got, exists := entryMap[key]; !exists {
			t.Errorf("Apply() missing key %q", key)
		} else if got != expected {
			t.Errorf("Apply() key %q = %q, want %q", key, got, expected)
		}
	}
}

func TestPatchouliParser_SupportedTypes(t *testing.T) {
	parser := NewPatchouliParser()

	types := parser.SupportedTypes()
	if len(types) == 0 {
		t.Error("SupportedTypes() returned empty slice")
	}

	expectedTypes := []string{"patchouli_category", "patchouli_entry"}
	for _, expected := range expectedTypes {
		found := false
		for _, typ := range types {
			if typ == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedTypes() does not include %q", expected)
		}
	}
}

func TestPatchouliParser_TagDetection(t *testing.T) {
	parser := NewPatchouliParser()

	tests := []struct {
		name        string
		content     string
		expectedTag string
	}{
		{
			name: "category_detection",
			content: `{
				"name": "Getting Started",
				"description": "Learn the basics"
			}`,
			expectedTag: "patchouli",
		},
		{
			name: "entry_detection",
			content: `{
				"name": "Welcome",
				"pages": [
					{"type": "text", "text": "Content"}
				]
			}`,
			expectedTag: "patchouli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := parser.Parse([]byte(tt.content))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(entries) == 0 {
				t.Fatal("Parse() returned no entries")
			}

			// Check that at least one entry has the expected tag
			foundTag := false
			for _, entry := range entries {
				for _, tag := range entry.Tags {
					if tag == tt.expectedTag {
						foundTag = true
						break
					}
				}
				if foundTag {
					break
				}
			}

			if !foundTag {
				t.Errorf("Expected tag %q not found in any entry", tt.expectedTag)
			}
		})
	}
}
