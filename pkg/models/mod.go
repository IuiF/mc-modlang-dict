// Package models defines the core data structures for the dictionary.
package models

import "time"

// Mod represents a Minecraft mod's metadata.
type Mod struct {
	ID          string            `json:"id" yaml:"id" gorm:"primaryKey"`
	DisplayName string            `json:"display_name" yaml:"display_name"`
	Author      string            `json:"author,omitempty" yaml:"author,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty" gorm:"serializer:json"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" gorm:"serializer:json"`
	CreatedAt   time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Mod) TableName() string {
	return "mods"
}
