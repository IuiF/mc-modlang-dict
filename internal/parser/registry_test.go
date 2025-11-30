package parser

import (
	"testing"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// mockParser is a simple parser implementation for testing.
type mockParser struct {
	name  string
	types []string
}

func (m *mockParser) Parse(content []byte) ([]interfaces.ParsedEntry, error) {
	return nil, nil
}

func (m *mockParser) Apply(content []byte, translations map[string]string) ([]byte, error) {
	return content, nil
}

func (m *mockParser) SupportedTypes() []string {
	return m.types
}

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Fatal("NewRegistry() returned nil")
	}
}

func TestRegistry_Register_Get(t *testing.T) {
	reg := NewRegistry()

	parser := &mockParser{name: "json_lang", types: []string{"lang"}}
	reg.Register("json_lang", parser)

	got, ok := reg.Get("json_lang")
	if !ok {
		t.Error("Get() returned false for registered parser")
	}
	if got != parser {
		t.Error("Get() returned different parser instance")
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	reg := NewRegistry()

	_, ok := reg.Get("nonexistent")
	if ok {
		t.Error("Get() returned true for unregistered parser")
	}
}

func TestRegistry_List(t *testing.T) {
	reg := NewRegistry()

	reg.Register("json_lang", &mockParser{})
	reg.Register("patchouli", &mockParser{})
	reg.Register("snbt", &mockParser{})

	names := reg.List()
	if len(names) != 3 {
		t.Errorf("List() returned %d names, want 3", len(names))
	}

	// Check all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	for _, expected := range []string{"json_lang", "patchouli", "snbt"} {
		if !nameMap[expected] {
			t.Errorf("List() missing %q", expected)
		}
	}
}

func TestRegistry_Register_Override(t *testing.T) {
	reg := NewRegistry()

	parser1 := &mockParser{name: "v1"}
	parser2 := &mockParser{name: "v2"}

	reg.Register("test", parser1)
	reg.Register("test", parser2)

	got, _ := reg.Get("test")
	if got != parser2 {
		t.Error("Register() did not override existing parser")
	}
}

func TestNewDefaultRegistry(t *testing.T) {
	reg := NewDefaultRegistry()

	// Check json_lang is registered
	if _, ok := reg.Get("json_lang"); !ok {
		t.Error("json_lang parser not registered")
	}

	// Check patchouli is registered
	if _, ok := reg.Get("patchouli"); !ok {
		t.Error("patchouli parser not registered")
	}

	// Check mantle_book is registered
	if _, ok := reg.Get("mantle_book"); !ok {
		t.Error("mantle_book parser not registered")
	}
}
