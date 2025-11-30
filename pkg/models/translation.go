package models

import "time"

// Translation represents a translated string.
// During migration period, both old (ModVersionID-based) and new (SourceID-based) fields are supported.
type Translation struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SourceID   int64     `json:"source_id" gorm:"index"`                   // New: Links to TranslationSource
	TargetText *string   `json:"target_text,omitempty"`
	TargetLang string    `json:"target_lang" gorm:"default:ja_jp;index"`
	Status     string    `json:"status" gorm:"default:pending;index"`      // pending, translated, verified, inherited, needs_review
	Translator *string   `json:"translator,omitempty"`                     // "claude", "community", "official"
	Tags       []string  `json:"tags,omitempty" gorm:"serializer:json"`
	Notes      *string   `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Legacy fields - used during migration and by translate.go
	ModVersionID int64  `json:"mod_version_id" gorm:"index"`
	Key          string `json:"key"`
	SourceText   string `json:"source_text"`
	SourceLang   string `json:"source_lang" gorm:"default:en_us"`
}

// TableName returns the table name for GORM.
func (Translation) TableName() string {
	return "translations"
}

// Translation status constants.
const (
	StatusPending     = "pending"
	StatusTranslated  = "translated"
	StatusVerified    = "verified"
	StatusInherited   = "inherited"    // Inherited from previous version (identical source)
	StatusNeedsReview = "needs_review" // Source text changed, translation may be outdated
)

// TranslationWithSource combines Translation with its source information.
// Used for queries that need both translation and source data.
type TranslationWithSource struct {
	Translation
	Key        string `json:"key"`
	SourceText string `json:"source_text"`
	SourceLang string `json:"source_lang"`
	IsCurrent  bool   `json:"is_current"`
}
