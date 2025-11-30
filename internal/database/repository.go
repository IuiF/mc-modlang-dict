// Package database provides SQLite implementation of the Repository interface.
package database

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// Errors
var (
	ErrModNotFound         = errors.New("mod not found")
	ErrTermNotFound        = errors.New("term not found")
	ErrVersionNotFound     = errors.New("version not found")
	ErrTranslationNotFound = errors.New("translation not found")
	ErrPatternNotFound     = errors.New("pattern not found")
)

// Repository implements interfaces.Repository using SQLite.
type Repository struct {
	db *gorm.DB
}

// Compile-time check that Repository implements interfaces.Repository.
var _ interfaces.Repository = (*Repository)(nil)

// NewRepository creates a new SQLite repository.
// Use ":memory:" for in-memory database or a file path for persistent storage.
func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Repository{db: db}, nil
}

// Close closes the database connection.
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Migrate runs database migrations.
func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(
		&models.Mod{},
		&models.ModVersion{},
		&models.Term{},
		&models.Translation{},
		&models.FilePattern{},
		&models.VersionDiff{},
	)
}

// GetMod retrieves a mod by ID.
func (r *Repository) GetMod(ctx context.Context, modID string) (*models.Mod, error) {
	var mod models.Mod
	if err := r.db.WithContext(ctx).First(&mod, "id = ?", modID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModNotFound
		}
		return nil, fmt.Errorf("failed to get mod %s: %w", modID, err)
	}
	return &mod, nil
}

// ListMods retrieves mods matching the filter.
func (r *Repository) ListMods(ctx context.Context, filter interfaces.ModFilter) ([]*models.Mod, error) {
	var mods []*models.Mod
	query := r.db.WithContext(ctx)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&mods).Error; err != nil {
		return nil, fmt.Errorf("failed to list mods: %w", err)
	}

	return mods, nil
}

// SaveMod creates or updates a mod.
func (r *Repository) SaveMod(ctx context.Context, mod *models.Mod) error {
	if err := r.db.WithContext(ctx).Save(mod).Error; err != nil {
		return fmt.Errorf("failed to save mod %s: %w", mod.ID, err)
	}
	return nil
}

// DeleteMod deletes a mod by ID.
func (r *Repository) DeleteMod(ctx context.Context, modID string) error {
	result := r.db.WithContext(ctx).Delete(&models.Mod{}, "id = ?", modID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete mod %s: %w", modID, result.Error)
	}
	return nil
}

// Version operations

// GetVersion retrieves a version by ID.
func (r *Repository) GetVersion(ctx context.Context, id int64) (*models.ModVersion, error) {
	var version models.ModVersion
	if err := r.db.WithContext(ctx).First(&version, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, fmt.Errorf("failed to get version %d: %w", id, err)
	}
	return &version, nil
}

// GetVersionBySpec retrieves a version by mod ID, version, and MC version.
func (r *Repository) GetVersionBySpec(ctx context.Context, modID, version, mcVersion string) (*models.ModVersion, error) {
	var v models.ModVersion
	err := r.db.WithContext(ctx).
		Where("mod_id = ? AND version = ? AND mc_version = ?", modID, version, mcVersion).
		First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return &v, nil
}

// ListVersions retrieves versions for a mod matching the filter.
func (r *Repository) ListVersions(ctx context.Context, modID string, filter interfaces.VersionFilter) ([]*models.ModVersion, error) {
	var versions []*models.ModVersion
	query := r.db.WithContext(ctx).Where("mod_id = ?", modID)

	if filter.MCVersion != "" {
		query = query.Where("mc_version = ?", filter.MCVersion)
	}
	if filter.Loader != "" {
		query = query.Where("loader = ?", filter.Loader)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	return versions, nil
}

// SaveVersion creates or updates a version.
func (r *Repository) SaveVersion(ctx context.Context, version *models.ModVersion) error {
	if err := r.db.WithContext(ctx).Save(version).Error; err != nil {
		return fmt.Errorf("failed to save version: %w", err)
	}
	return nil
}

// DeleteVersion deletes a version by ID.
func (r *Repository) DeleteVersion(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.ModVersion{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete version %d: %w", id, result.Error)
	}
	return nil
}

// Term operations

// GetTerm retrieves a term by ID.
func (r *Repository) GetTerm(ctx context.Context, id int64) (*models.Term, error) {
	var term models.Term
	if err := r.db.WithContext(ctx).First(&term, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTermNotFound
		}
		return nil, fmt.Errorf("failed to get term %d: %w", id, err)
	}
	return &term, nil
}

// ListTerms retrieves terms matching the filter.
func (r *Repository) ListTerms(ctx context.Context, filter interfaces.TermFilter) ([]*models.Term, error) {
	var terms []*models.Term
	query := r.db.WithContext(ctx)

	// Filter by single scope
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}

	// Filter by multiple scopes
	if len(filter.Scopes) > 0 {
		query = query.Where("scope IN ?", filter.Scopes)
	}

	// Filter by target language
	if filter.TargetLang != "" {
		query = query.Where("target_lang = ?", filter.TargetLang)
	}

	// Pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by priority (higher first)
	query = query.Order("priority DESC")

	if err := query.Find(&terms).Error; err != nil {
		return nil, fmt.Errorf("failed to list terms: %w", err)
	}

	return terms, nil
}

// SaveTerm creates or updates a term.
func (r *Repository) SaveTerm(ctx context.Context, term *models.Term) error {
	if err := r.db.WithContext(ctx).Save(term).Error; err != nil {
		return fmt.Errorf("failed to save term: %w", err)
	}
	return nil
}

// DeleteTerm deletes a term by ID.
func (r *Repository) DeleteTerm(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Term{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete term %d: %w", id, result.Error)
	}
	return nil
}

// BulkSaveTerms creates or updates multiple terms in a transaction.
func (r *Repository) BulkSaveTerms(ctx context.Context, terms []*models.Term) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, term := range terms {
			if err := tx.Save(term).Error; err != nil {
				return fmt.Errorf("failed to save term: %w", err)
			}
		}
		return nil
	})
}

// Translation operations

// GetTranslation retrieves a translation by version ID and key.
func (r *Repository) GetTranslation(ctx context.Context, versionID int64, key string) (*models.Translation, error) {
	var trans models.Translation
	err := r.db.WithContext(ctx).
		Where("mod_version_id = ? AND key = ?", versionID, key).
		First(&trans).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTranslationNotFound
		}
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}
	return &trans, nil
}

// ListTranslations retrieves translations for a version matching the filter.
func (r *Repository) ListTranslations(ctx context.Context, versionID int64, filter interfaces.TranslationFilter) ([]*models.Translation, error) {
	var translations []*models.Translation
	query := r.db.WithContext(ctx).Where("mod_version_id = ?", versionID)

	if filter.TargetLang != "" {
		query = query.Where("target_lang = ?", filter.TargetLang)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&translations).Error; err != nil {
		return nil, fmt.Errorf("failed to list translations: %w", err)
	}

	return translations, nil
}

// SaveTranslation creates or updates a translation.
func (r *Repository) SaveTranslation(ctx context.Context, translation *models.Translation) error {
	if err := r.db.WithContext(ctx).Save(translation).Error; err != nil {
		return fmt.Errorf("failed to save translation: %w", err)
	}
	return nil
}

// BulkSaveTranslations creates or updates multiple translations in a transaction.
func (r *Repository) BulkSaveTranslations(ctx context.Context, translations []*models.Translation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, trans := range translations {
			if err := tx.Save(trans).Error; err != nil {
				return fmt.Errorf("failed to save translation: %w", err)
			}
		}
		return nil
	})
}

// DeleteTranslation deletes a translation by ID.
func (r *Repository) DeleteTranslation(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Translation{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete translation %d: %w", id, result.Error)
	}
	return nil
}

// Pattern operations

// GetPattern retrieves a pattern by ID.
func (r *Repository) GetPattern(ctx context.Context, id int64) (*models.FilePattern, error) {
	var pattern models.FilePattern
	if err := r.db.WithContext(ctx).First(&pattern, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatternNotFound
		}
		return nil, fmt.Errorf("failed to get pattern %d: %w", id, err)
	}
	return &pattern, nil
}

// ListPatterns retrieves patterns for a scope.
func (r *Repository) ListPatterns(ctx context.Context, scope string) ([]*models.FilePattern, error) {
	var patterns []*models.FilePattern
	query := r.db.WithContext(ctx)

	if scope != "" {
		query = query.Where("scope = ?", scope)
	}

	query = query.Order("priority DESC")

	if err := query.Find(&patterns).Error; err != nil {
		return nil, fmt.Errorf("failed to list patterns: %w", err)
	}

	return patterns, nil
}

// SavePattern creates or updates a pattern.
func (r *Repository) SavePattern(ctx context.Context, pattern *models.FilePattern) error {
	if err := r.db.WithContext(ctx).Save(pattern).Error; err != nil {
		return fmt.Errorf("failed to save pattern: %w", err)
	}
	return nil
}

// DeletePattern deletes a pattern by ID.
func (r *Repository) DeletePattern(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.FilePattern{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete pattern %d: %w", id, result.Error)
	}
	return nil
}

// Diff operations

// ListDiffs retrieves diffs between two versions.
func (r *Repository) ListDiffs(ctx context.Context, fromVersionID, toVersionID int64) ([]*models.VersionDiff, error) {
	var diffs []*models.VersionDiff
	err := r.db.WithContext(ctx).
		Where("from_version_id = ? AND to_version_id = ?", fromVersionID, toVersionID).
		Find(&diffs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list diffs: %w", err)
	}
	return diffs, nil
}

// SaveDiff creates or updates a diff.
func (r *Repository) SaveDiff(ctx context.Context, diff *models.VersionDiff) error {
	if err := r.db.WithContext(ctx).Save(diff).Error; err != nil {
		return fmt.Errorf("failed to save diff: %w", err)
	}
	return nil
}

// BulkSaveDiffs creates or updates multiple diffs in a transaction.
func (r *Repository) BulkSaveDiffs(ctx context.Context, diffs []*models.VersionDiff) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, diff := range diffs {
			if err := tx.Save(diff).Error; err != nil {
				return fmt.Errorf("failed to save diff: %w", err)
			}
		}
		return nil
	})
}
