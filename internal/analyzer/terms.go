package analyzer

import (
	"context"
	"sort"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// checkTermCompliance checks if translations comply with term dictionary.
func (a *Analyzer) checkTermCompliance(ctx context.Context, modID string) ([]TermViolation, error) {
	// Get all applicable terms
	terms, err := a.getApplicableTerms(ctx, modID)
	if err != nil {
		return nil, err
	}

	if len(terms) == 0 {
		return nil, nil
	}

	// Get all translation pairs
	pairs, err := a.getAllTranslationPairs(ctx, modID)
	if err != nil {
		return nil, err
	}

	// Check each translation against terms
	violationMap := make(map[string]*TermViolation)

	for _, pair := range pairs {
		for _, term := range terms {
			// Check if source_text contains the term
			if !containsWord(pair.SourceText, term.SourceText) {
				continue
			}

			// Check if target_text contains the expected translation
			if containsJapanese(pair.TargetText, term.TargetText) {
				continue // Compliant
			}

			// This is a violation
			key := term.SourceText + ":" + term.TargetText
			if violationMap[key] == nil {
				violationMap[key] = &TermViolation{
					TermSource: term.SourceText,
					TermTarget: term.TargetText,
					TermScope:  term.Scope,
					Violations: []ViolationDetail{},
				}
			}

			// Add the violation detail
			violationMap[key].Violations = append(violationMap[key].Violations, ViolationDetail{
				Key:        pair.Key,
				SourceText: pair.SourceText,
				TargetText: pair.TargetText,
				Expected:   term.TargetText,
			})
			violationMap[key].ViolationCount++
		}
	}

	// Convert to list and sort by violation count
	var violations []TermViolation
	for _, v := range violationMap {
		violations = append(violations, *v)
	}

	sort.Slice(violations, func(i, j int) bool {
		return violations[i].ViolationCount > violations[j].ViolationCount
	})

	return violations, nil
}

// getApplicableTerms returns terms applicable to the given mod.
func (a *Analyzer) getApplicableTerms(ctx context.Context, modID string) ([]*models.Term, error) {
	// Build filter for applicable scopes
	filter := interfaces.TermFilter{}

	// Get global and mod-specific terms
	allTerms, err := a.repo.ListTerms(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Filter to applicable scopes
	var applicable []*models.Term
	for _, term := range allTerms {
		switch {
		case term.Scope == "global":
			applicable = append(applicable, term)
		case modID != "" && term.Scope == "mod:"+modID:
			applicable = append(applicable, term)
		case strings.HasPrefix(term.Scope, "category:"):
			// Category terms are always applicable for now
			// TODO: Check mod tags to determine category membership
			applicable = append(applicable, term)
		}
	}

	// Sort by priority (higher first)
	sort.Slice(applicable, func(i, j int) bool {
		return applicable[i].Priority > applicable[j].Priority
	})

	return applicable, nil
}

// containsWord checks if text contains the word (case-insensitive, word boundary aware).
func containsWord(text, word string) bool {
	textLower := strings.ToLower(text)
	wordLower := strings.ToLower(word)

	// Simple contains check
	idx := strings.Index(textLower, wordLower)
	if idx == -1 {
		return false
	}

	// Check word boundaries
	// Start boundary
	if idx > 0 {
		prevChar := text[idx-1]
		if isWordChar(prevChar) {
			// Not at word boundary, search for another occurrence
			return containsWord(text[idx+1:], word)
		}
	}

	// End boundary
	endIdx := idx + len(word)
	if endIdx < len(text) {
		nextChar := text[endIdx]
		if isWordChar(nextChar) {
			// Not at word boundary
			return containsWord(text[idx+1:], word)
		}
	}

	return true
}

// isWordChar checks if a byte is a word character.
func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// containsJapanese checks if text contains the Japanese string.
func containsJapanese(text, substr string) bool {
	return strings.Contains(text, substr)
}
