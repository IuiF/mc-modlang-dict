package models

// FilePattern defines a pattern for locating translation files in a mod.
type FilePattern struct {
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	Scope       string  `json:"scope" gorm:"index;not null"` // "global", "mod:{mod_id}"
	Pattern     string  `json:"pattern" gorm:"not null"`     // e.g., "assets/{mod_id}/lang/{lang}.json"
	Type        string  `json:"type" gorm:"not null"`        // "lang", "book", "manual", "quest"
	Parser      string  `json:"parser" gorm:"not null"`      // "json_lang", "patchouli", "snbt"
	Priority    int     `json:"priority" gorm:"default:100"`
	Required    bool    `json:"required" gorm:"default:false"`
	Description *string `json:"description,omitempty"`
}

// TableName returns the table name for GORM.
func (FilePattern) TableName() string {
	return "file_patterns"
}

// Pattern type constants.
const (
	PatternTypeLang   = "lang"
	PatternTypeBook   = "book"
	PatternTypeManual = "manual"
	PatternTypeQuest  = "quest"
	PatternTypeData   = "data"
)

// Parser type constants.
const (
	ParserJSONLang    = "json_lang"
	ParserJSONGeneric = "json_generic"
	ParserPatchouli   = "patchouli"
	ParserSNBT        = "snbt"
	ParserLegacyLang  = "lang_legacy"
)
