package models

import "time"

// VersionDiff represents a difference between two mod versions.
type VersionDiff struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	FromVersionID int64     `json:"from_version_id" gorm:"index;not null"`
	ToVersionID   int64     `json:"to_version_id" gorm:"index;not null"`
	Type          string    `json:"type" gorm:"not null"` // added, removed, changed
	Key           string    `json:"key" gorm:"not null"`
	OldText       *string   `json:"old_text,omitempty"`
	NewText       *string   `json:"new_text,omitempty"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for GORM.
func (VersionDiff) TableName() string {
	return "version_diffs"
}

// Diff type constants.
const (
	DiffTypeAdded   = "added"
	DiffTypeRemoved = "removed"
	DiffTypeChanged = "changed"
)
