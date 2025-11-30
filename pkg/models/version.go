package models

import "time"

// ModVersion represents a specific version of a mod.
type ModVersion struct {
	ID              int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	ModID           string            `json:"mod_id" gorm:"index;not null"`
	Version         string            `json:"version" gorm:"not null"`
	MCVersion       string            `json:"mc_version" gorm:"not null"`
	Loader          string            `json:"loader"`                              // forge, fabric, neoforge, quilt
	IsDefault       bool              `json:"is_default" gorm:"default:false"`     // Default version for this mod
	ParentVersionID *int64            `json:"parent_version_id,omitempty"`
	Stats           VersionStats      `json:"stats" gorm:"embedded"`
	Metadata        map[string]string `json:"metadata,omitempty" gorm:"serializer:json"`
	CreatedAt       time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// VersionStats holds statistics about a mod version's translations.
type VersionStats struct {
	TotalKeys      int `json:"total_keys"`
	TranslatedKeys int `json:"translated_keys"`
	VerifiedKeys   int `json:"verified_keys"`
}

// TableName returns the table name for GORM.
func (ModVersion) TableName() string {
	return "mod_versions"
}
