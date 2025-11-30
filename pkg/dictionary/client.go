// Package dictionary provides the main API for the mod dictionary.
package dictionary

import (
	"context"
	"fmt"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// Client is the main entry point for dictionary operations.
type Client struct {
	repo    interfaces.Repository
	parsers interfaces.ParserRegistry
	config  *Config
}

// New creates a new dictionary client with the given options.
func New(repo interfaces.Repository, opts ...Option) (*Client, error) {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	return &Client{
		repo:   repo,
		config: config,
	}, nil
}

// Close closes the dictionary client and releases resources.
func (c *Client) Close() error {
	if c.repo != nil {
		return c.repo.Close()
	}
	return nil
}

// ==================== Mod Operations ====================

// GetMod retrieves a mod by ID.
func (c *Client) GetMod(ctx context.Context, modID string) (*models.Mod, error) {
	return c.repo.GetMod(ctx, modID)
}

// ListMods lists all registered mods.
func (c *Client) ListMods(ctx context.Context, filter interfaces.ModFilter) ([]*models.Mod, error) {
	return c.repo.ListMods(ctx, filter)
}

// SaveMod saves or updates a mod.
func (c *Client) SaveMod(ctx context.Context, mod *models.Mod) error {
	return c.repo.SaveMod(ctx, mod)
}

// ==================== Version Operations ====================

// GetVersion retrieves a mod version by ID.
func (c *Client) GetVersion(ctx context.Context, id int64) (*models.ModVersion, error) {
	return c.repo.GetVersion(ctx, id)
}

// GetVersionBySpec retrieves a mod version by specification.
func (c *Client) GetVersionBySpec(ctx context.Context, modID, version, mcVersion string) (*models.ModVersion, error) {
	return c.repo.GetVersionBySpec(ctx, modID, version, mcVersion)
}

// ListVersions lists all versions for a mod.
func (c *Client) ListVersions(ctx context.Context, modID string, filter interfaces.VersionFilter) ([]*models.ModVersion, error) {
	return c.repo.ListVersions(ctx, modID, filter)
}

// ==================== Term Operations ====================

// TermQuery defines parameters for term queries.
type TermQuery struct {
	ModID         *string  // If set, includes mod-specific terms
	Categories    []string // Category scopes to include
	TargetLang    string   // Required: target language
	Tags          []string // Optional tag filter
	IncludeGlobal bool     // Include global terms (default: true)
}

// GetTerms retrieves terms based on the query.
// Terms are returned in priority order (global < category < mod).
func (c *Client) GetTerms(ctx context.Context, q TermQuery) ([]*models.Term, error) {
	var scopes []string

	// Build scope list in priority order (lowest first)
	if q.IncludeGlobal {
		scopes = append(scopes, models.ScopeGlobal)
	}
	for _, cat := range q.Categories {
		scopes = append(scopes, models.BuildScope(models.ScopeCategory, cat))
	}
	if q.ModID != nil && *q.ModID != "" {
		scopes = append(scopes, models.BuildScope(models.ScopeMod, *q.ModID))
	}

	return c.repo.ListTerms(ctx, interfaces.TermFilter{
		Scopes:     scopes,
		TargetLang: q.TargetLang,
		Tags:       q.Tags,
	})
}

// FormatTermsForLLM formats terms for use in LLM prompts.
func (c *Client) FormatTermsForLLM(ctx context.Context, q TermQuery) (string, error) {
	terms, err := c.GetTerms(ctx, q)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, term := range terms {
		sb.WriteString(fmt.Sprintf("- %s → %s", term.SourceText, term.TargetText))
		if term.Context != nil && *term.Context != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", *term.Context))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// ==================== Translation Operations ====================

// GetTranslation retrieves a specific translation.
func (c *Client) GetTranslation(ctx context.Context, versionID int64, key string) (*models.Translation, error) {
	return c.repo.GetTranslation(ctx, versionID, key)
}

// ListTranslations retrieves translations for a mod version.
func (c *Client) ListTranslations(ctx context.Context, versionID int64, filter interfaces.TranslationFilter) ([]*models.Translation, error) {
	return c.repo.ListTranslations(ctx, versionID, filter)
}

// BulkGetTranslations retrieves multiple translations efficiently.
func (c *Client) BulkGetTranslations(ctx context.Context, versionID int64, keys []string) (map[string]*models.Translation, error) {
	translations, err := c.repo.ListTranslations(ctx, versionID, interfaces.TranslationFilter{})
	if err != nil {
		return nil, err
	}

	// Build key set for filtering
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	result := make(map[string]*models.Translation)
	for _, t := range translations {
		if _, ok := keySet[t.Key]; ok {
			result[t.Key] = t
		}
	}

	return result, nil
}

// ==================== Pattern Operations ====================

// GetPatterns retrieves file patterns for a mod.
func (c *Client) GetPatterns(ctx context.Context, modID string) ([]*models.FilePattern, error) {
	// Get global patterns
	global, err := c.repo.ListPatterns(ctx, models.ScopeGlobal)
	if err != nil {
		return nil, err
	}

	// Get mod-specific patterns
	modScope := models.BuildScope(models.ScopeMod, modID)
	modPatterns, err := c.repo.ListPatterns(ctx, modScope)
	if err != nil {
		return nil, err
	}

	// Merge: mod patterns override global
	return append(global, modPatterns...), nil
}

// ==================== Statistics ====================

// Stats holds dictionary statistics.
type Stats struct {
	ModCount         int
	VersionCount     int
	TermCount        int
	TranslationCount int
	GlobalTerms      int
	CategoryTerms    int
	ModTerms         int
}

// Stats returns dictionary statistics.
func (c *Client) Stats(ctx context.Context) (*Stats, error) {
	// Implementation would query counts from repository
	// Placeholder for now
	return &Stats{}, nil
}
