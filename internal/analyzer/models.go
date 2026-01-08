// Package analyzer provides translation consistency analysis functionality.
package analyzer

import "time"

// AnalysisResult is the complete result of all analysis types.
type AnalysisResult struct {
	AnalysisDate  time.Time       `json:"analysis_date"`
	TargetMod     string          `json:"target_mod,omitempty"`
	Summary       AnalysisSummary `json:"summary"`
	Consistency   []ConsistencyIssue  `json:"consistency,omitempty"`
	Phrases       []PhraseIssue       `json:"discovered_phrases,omitempty"`
	TermViolations []TermViolation    `json:"term_violations,omitempty"`
}

// AnalysisSummary provides overview statistics.
type AnalysisSummary struct {
	TotalTranslations    int `json:"total_translations"`
	ConsistencyIssues    int `json:"consistency_issues"`
	DiscoveredPhrases    int `json:"discovered_phrases"`
	InconsistentPhrases  int `json:"inconsistent_phrases"`
	TermViolations       int `json:"term_violations"`
}

// ConsistencyIssue represents a case where the same source text has multiple translations.
type ConsistencyIssue struct {
	ModID        string            `json:"mod_id"`
	SourceText   string            `json:"source_text"`
	Translations []string          `json:"translations"`
	Counts       map[string]int    `json:"counts"`
	Keys         []string          `json:"keys,omitempty"`
	Suggested    string            `json:"suggested"`
	Reason       string            `json:"reason,omitempty"`
}

// PhraseIssue represents a discovered phrase pattern with its translations.
type PhraseIssue struct {
	Phrase           string            `json:"phrase"`
	Occurrences      int               `json:"occurrences"`
	Translations     map[string]int    `json:"translations"`
	ConsistencyScore float64           `json:"consistency_score"`
	InDictionary     bool              `json:"in_dictionary"`
	Suggested        string            `json:"suggested"`
	Examples         []PhraseExample   `json:"examples,omitempty"`
}

// PhraseExample shows a specific usage of a phrase.
type PhraseExample struct {
	Key        string `json:"key,omitempty"`
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
}

// TermViolation represents a case where a translation violates term dictionary.
type TermViolation struct {
	TermSource   string            `json:"term_source"`
	TermTarget   string            `json:"term_target"`
	TermScope    string            `json:"term_scope"`
	Violations   []ViolationDetail `json:"violations"`
	ViolationCount int             `json:"violation_count"`
}

// ViolationDetail provides details about a specific term violation.
type ViolationDetail struct {
	Key        string `json:"key"`
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
	Expected   string `json:"expected"`
}

// SourceTranslationPair represents a source text with its translation.
type SourceTranslationPair struct {
	ModID      string
	Key        string
	SourceText string
	TargetText string
	Status     string
}

// AnalysisOptions contains options for analysis.
type AnalysisOptions struct {
	ModID      string
	MinCount   int
	Format     string
	OutputPath string
}
