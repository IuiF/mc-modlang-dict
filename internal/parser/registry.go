// Package parser provides translation file parsers and a registry for managing them.
package parser

import (
	"sort"
	"sync"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// Registry implements interfaces.ParserRegistry.
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]interfaces.Parser
}

// Compile-time check that Registry implements interfaces.ParserRegistry.
var _ interfaces.ParserRegistry = (*Registry)(nil)

// NewRegistry creates a new parser registry.
func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]interfaces.Parser),
	}
}

// Register adds a parser to the registry.
func (r *Registry) Register(name string, parser interfaces.Parser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[name] = parser
}

// Get retrieves a parser by name.
func (r *Registry) Get(name string) (interfaces.Parser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parser, ok := r.parsers[name]
	return parser, ok
}

// List returns all registered parser names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.parsers))
	for name := range r.parsers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// NewDefaultRegistry creates a registry with all built-in parsers registered.
func NewDefaultRegistry() *Registry {
	reg := NewRegistry()
	reg.Register("json_lang", NewJSONLangParser())
	reg.Register("patchouli", NewPatchouliParser())
	reg.Register("mantle_book", NewMantleBookParser())
	return reg
}
