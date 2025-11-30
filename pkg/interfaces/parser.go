package interfaces

// Parser defines the interface for translation file parsers.
type Parser interface {
	// Parse extracts translation entries from file content.
	Parse(content []byte) ([]ParsedEntry, error)

	// Apply applies translations to the original content.
	Apply(content []byte, translations map[string]string) ([]byte, error)

	// SupportedTypes returns the pattern types this parser handles.
	SupportedTypes() []string
}

// ParsedEntry represents a single translatable entry extracted from a file.
type ParsedEntry struct {
	Key        string   // Translation key
	Text       string   // Original text
	Context    *string  // Optional context information
	Tags       []string // Auto-detected tags
	LineNumber int      // Source line number
	FilePath   string   // Source file path (for multi-file parsing)
}

// ParserRegistry manages available parsers.
type ParserRegistry interface {
	// Register adds a parser to the registry.
	Register(name string, parser Parser)

	// Get retrieves a parser by name.
	Get(name string) (Parser, bool)

	// List returns all registered parser names.
	List() []string
}
