package models

import "time"

// Term represents a translation term in the dictionary.
type Term struct {
	ID         int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Scope      string   `json:"scope" gorm:"index;not null"` // "global", "category:{name}", "mod:{mod_id}"
	SourceText string   `json:"source_text" gorm:"not null"`
	TargetText string   `json:"target_text" gorm:"not null"`
	SourceLang string   `json:"source_lang" gorm:"default:en_us"`
	TargetLang string   `json:"target_lang" gorm:"default:ja_jp;index"`
	Context    *string  `json:"context,omitempty"`
	Tags       []string `json:"tags,omitempty" gorm:"serializer:json"`
	Priority   int      `json:"priority" gorm:"default:100"`
	Source     string   `json:"source"` // "official", "community", "claude"
	Notes      *string  `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for GORM.
func (Term) TableName() string {
	return "terms"
}

// ScopeType constants for term scope.
const (
	ScopeGlobal   = "global"
	ScopeCategory = "category"
	ScopeMod      = "mod"
)

// ParseScope extracts the scope type and value from a scope string.
// Examples: "global" -> ("global", ""), "category:tech" -> ("category", "tech")
func ParseScope(scope string) (scopeType, value string) {
	if scope == ScopeGlobal {
		return ScopeGlobal, ""
	}
	// Find colon separator
	for i, c := range scope {
		if c == ':' {
			return scope[:i], scope[i+1:]
		}
	}
	return scope, ""
}

// BuildScope creates a scope string from type and value.
func BuildScope(scopeType, value string) string {
	if scopeType == ScopeGlobal || value == "" {
		return ScopeGlobal
	}
	return scopeType + ":" + value
}
