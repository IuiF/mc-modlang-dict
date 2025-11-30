// Package interfaces defines the contracts for data access and parsing.
package interfaces

import (
	"context"

	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// Repository defines the data access interface.
type Repository interface {
	// Mod operations
	GetMod(ctx context.Context, modID string) (*models.Mod, error)
	ListMods(ctx context.Context, filter ModFilter) ([]*models.Mod, error)
	SaveMod(ctx context.Context, mod *models.Mod) error
	DeleteMod(ctx context.Context, modID string) error

	// Version operations
	GetVersion(ctx context.Context, id int64) (*models.ModVersion, error)
	GetVersionBySpec(ctx context.Context, modID, version, mcVersion string) (*models.ModVersion, error)
	ListVersions(ctx context.Context, modID string, filter VersionFilter) ([]*models.ModVersion, error)
	SaveVersion(ctx context.Context, version *models.ModVersion) error
	DeleteVersion(ctx context.Context, id int64) error

	// Term operations
	GetTerm(ctx context.Context, id int64) (*models.Term, error)
	ListTerms(ctx context.Context, filter TermFilter) ([]*models.Term, error)
	SaveTerm(ctx context.Context, term *models.Term) error
	DeleteTerm(ctx context.Context, id int64) error
	BulkSaveTerms(ctx context.Context, terms []*models.Term) error

	// Translation operations
	GetTranslation(ctx context.Context, versionID int64, key string) (*models.Translation, error)
	ListTranslations(ctx context.Context, versionID int64, filter TranslationFilter) ([]*models.Translation, error)
	SaveTranslation(ctx context.Context, translation *models.Translation) error
	BulkSaveTranslations(ctx context.Context, translations []*models.Translation) error
	DeleteTranslation(ctx context.Context, id int64) error

	// Pattern operations
	GetPattern(ctx context.Context, id int64) (*models.FilePattern, error)
	ListPatterns(ctx context.Context, scope string) ([]*models.FilePattern, error)
	SavePattern(ctx context.Context, pattern *models.FilePattern) error
	DeletePattern(ctx context.Context, id int64) error

	// Diff operations
	ListDiffs(ctx context.Context, fromVersionID, toVersionID int64) ([]*models.VersionDiff, error)
	SaveDiff(ctx context.Context, diff *models.VersionDiff) error
	BulkSaveDiffs(ctx context.Context, diffs []*models.VersionDiff) error

	// Lifecycle
	Close() error
	Migrate() error
}

// ModFilter defines filter options for mod queries.
type ModFilter struct {
	Tags   []string
	Limit  int
	Offset int
}

// VersionFilter defines filter options for version queries.
type VersionFilter struct {
	MCVersion string
	Loader    string
	Limit     int
	Offset    int
}

// TermFilter defines filter options for term queries.
type TermFilter struct {
	Scope      string   // "global", "category:tech", "mod:create"
	Scopes     []string // Multiple scopes (merged results)
	TargetLang string
	Tags       []string
	Limit      int
	Offset     int
}

// TranslationFilter defines filter options for translation queries.
type TranslationFilter struct {
	TargetLang string
	Status     string
	Tags       []string
	Limit      int
	Offset     int
}
