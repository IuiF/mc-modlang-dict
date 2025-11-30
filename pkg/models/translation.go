package models

import "time"

// Translation represents a translated string for a specific mod version.
type Translation struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ModVersionID int64     `json:"mod_version_id" gorm:"index;not null"`
	Key          string    `json:"key" gorm:"not null"`          // Translation key (e.g., "item.modid.name")
	SourceText   string    `json:"source_text" gorm:"not null"`
	TargetText   *string   `json:"target_text,omitempty"`
	SourceLang   string    `json:"source_lang" gorm:"default:en_us"`
	TargetLang   string    `json:"target_lang" gorm:"default:ja_jp;index"`
	Status       string    `json:"status" gorm:"default:pending;index"` // pending, translated, verified
	Translator   *string   `json:"translator,omitempty"`                // "claude", "community", "official"
	Tags         []string  `json:"tags,omitempty" gorm:"serializer:json"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Translation) TableName() string {
	return "translations"
}

// Translation status constants.
const (
	StatusPending    = "pending"
	StatusTranslated = "translated"
	StatusVerified   = "verified"
	StatusInherited  = "inherited"   // Inherited from previous version (identical source)
	StatusNeedsReview = "needs_review" // Inherited but source text changed
)
