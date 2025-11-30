package dictionary

import "errors"

// Common errors.
var (
	ErrModNotFound         = errors.New("mod not found")
	ErrVersionNotFound     = errors.New("version not found")
	ErrTranslationNotFound = errors.New("translation not found")
	ErrTermNotFound        = errors.New("term not found")
	ErrPatternNotFound     = errors.New("pattern not found")
	ErrInvalidScope        = errors.New("invalid scope format")
	ErrInvalidParser       = errors.New("unknown parser type")
	ErrParseFailure        = errors.New("failed to parse file")
	ErrDatabaseConnection  = errors.New("database connection failed")
)
