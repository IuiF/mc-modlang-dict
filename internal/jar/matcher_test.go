package jar

import (
	"testing"
)

func TestNewMatcher(t *testing.T) {
	m := NewMatcher()
	if m == nil {
		t.Fatal("NewMatcher() returned nil")
	}
}

func TestMatcher_Match(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		path     string
		wantMatch bool
		wantVars map[string]string
	}{
		{
			name:     "basic_lang_file",
			pattern:  "assets/{mod_id}/lang/{lang}.json",
			path:     "assets/create/lang/en_us.json",
			wantMatch: true,
			wantVars: map[string]string{"mod_id": "create", "lang": "en_us"},
		},
		{
			name:     "different_mod",
			pattern:  "assets/{mod_id}/lang/{lang}.json",
			path:     "assets/botania/lang/ja_jp.json",
			wantMatch: true,
			wantVars: map[string]string{"mod_id": "botania", "lang": "ja_jp"},
		},
		{
			name:     "no_match_wrong_extension",
			pattern:  "assets/{mod_id}/lang/{lang}.json",
			path:     "assets/create/lang/en_us.txt",
			wantMatch: false,
			wantVars: nil,
		},
		{
			name:     "no_match_wrong_structure",
			pattern:  "assets/{mod_id}/lang/{lang}.json",
			path:     "assets/create/textures/item.png",
			wantMatch: false,
			wantVars: nil,
		},
		{
			name:     "patchouli_book",
			pattern:  "data/{mod_id}/patchouli_books/{book_id}/en_us/entries/{entry}.json",
			path:     "data/botania/patchouli_books/lexicon/en_us/entries/basics.json",
			wantMatch: true,
			wantVars: map[string]string{
				"mod_id":  "botania",
				"book_id": "lexicon",
				"entry":   "basics",
			},
		},
		{
			name:     "wildcard_pattern",
			pattern:  "assets/{mod_id}/lang/*.json",
			path:     "assets/create/lang/en_us.json",
			wantMatch: true,
			wantVars: map[string]string{"mod_id": "create"},
		},
		{
			name:     "double_wildcard",
			pattern:  "data/{mod_id}/**/*.json",
			path:     "data/botania/patchouli_books/lexicon/en_us/entries/basics.json",
			wantMatch: true,
			wantVars: map[string]string{"mod_id": "botania"},
		},
	}

	m := NewMatcher()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, vars := m.Match(tt.pattern, tt.path)

			if match != tt.wantMatch {
				t.Errorf("Match() = %v, want %v", match, tt.wantMatch)
			}

			if tt.wantMatch && tt.wantVars != nil {
				for key, want := range tt.wantVars {
					if got := vars[key]; got != want {
						t.Errorf("Match() vars[%q] = %q, want %q", key, got, want)
					}
				}
			}
		})
	}
}

func TestMatcher_FindFiles(t *testing.T) {
	m := NewMatcher()

	files := []string{
		"assets/create/lang/en_us.json",
		"assets/create/lang/ja_jp.json",
		"assets/create/textures/item.png",
		"assets/botania/lang/en_us.json",
		"META-INF/mods.toml",
	}

	pattern := "assets/{mod_id}/lang/{lang}.json"

	matches := m.FindFiles(pattern, files)

	if len(matches) != 3 {
		t.Errorf("FindFiles() found %d matches, want 3", len(matches))
	}

	// Check that all lang files are found
	expectedPaths := map[string]bool{
		"assets/create/lang/en_us.json":  true,
		"assets/create/lang/ja_jp.json":  true,
		"assets/botania/lang/en_us.json": true,
	}

	for _, match := range matches {
		if !expectedPaths[match.Path] {
			t.Errorf("FindFiles() unexpected match: %q", match.Path)
		}
	}
}

func TestMatcher_ExpandPattern(t *testing.T) {
	m := NewMatcher()

	tests := []struct {
		name    string
		pattern string
		vars    map[string]string
		want    string
	}{
		{
			name:    "basic",
			pattern: "assets/{mod_id}/lang/{lang}.json",
			vars:    map[string]string{"mod_id": "create", "lang": "ja_jp"},
			want:    "assets/create/lang/ja_jp.json",
		},
		{
			name:    "partial",
			pattern: "assets/{mod_id}/lang/{lang}.json",
			vars:    map[string]string{"mod_id": "create"},
			want:    "assets/create/lang/{lang}.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.ExpandPattern(tt.pattern, tt.vars)
			if got != tt.want {
				t.Errorf("ExpandPattern() = %q, want %q", got, tt.want)
			}
		})
	}
}
