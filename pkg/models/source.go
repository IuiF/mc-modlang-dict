package models

import "time"

// TranslationSource represents a unique source text for a translation key.
// Multiple versions can share the same source if the text hasn't changed.
type TranslationSource struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ModID      string    `json:"mod_id" gorm:"index;not null"`
	Key        string    `json:"key" gorm:"not null"`
	SourceText string    `json:"source_text" gorm:"not null"`
	SourceLang string    `json:"source_lang" gorm:"default:en_us"`
	IsCurrent  bool      `json:"is_current" gorm:"default:true;index"` // Current/default source for this key
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (TranslationSource) TableName() string {
	return "translation_sources"
}

// SourceVersion links a TranslationSource to a ModVersion (N:M relationship).
// This tracks which source text was used in which version.
type SourceVersion struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SourceID     int64     `json:"source_id" gorm:"index;not null"`
	ModVersionID int64     `json:"mod_version_id" gorm:"index;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for GORM.
func (SourceVersion) TableName() string {
	return "source_versions"
}
