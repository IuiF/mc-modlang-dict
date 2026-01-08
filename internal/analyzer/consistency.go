package analyzer

import (
	"context"
	"sort"
	"strings"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
	"github.com/iuif/minecraft-mod-dictionary/pkg/models"
)

// findConsistencyIssues finds cases where the same source text has different translations.
func (a *Analyzer) findConsistencyIssues(ctx context.Context, modID string) ([]ConsistencyIssue, int, error) {
	db := a.repo.GetDB()

	// Build the base query for getting all translation pairs
	type TranslationPair struct {
		ModID      string
		Key        string
		SourceText string
		TargetText string
		Status     string
	}

	var pairs []TranslationPair
	var totalCount int

	query := db.Table("translations").
		Select(`
			translation_sources.mod_id,
			translation_sources.key,
			translation_sources.source_text,
			translations.target_text,
			translations.status
		`).
		Joins("JOIN translation_sources ON translation_sources.id = translations.source_id").
		Joins("JOIN source_versions ON source_versions.source_id = translation_sources.id").
		Joins("JOIN mod_versions ON mod_versions.id = source_versions.mod_version_id").
		Where("mod_versions.is_default = ?", true).
		Where("translations.status IN ?", []string{"translated", "verified", "official"}).
		Where("translations.target_text IS NOT NULL").
		Where("translations.target_text != ''")

	if modID != "" {
		query = query.Where("translation_sources.mod_id = ?", modID)
	}

	if err := query.Scan(&pairs).Error; err != nil {
		return nil, 0, err
	}

	totalCount = len(pairs)

	// Group by (mod_id, source_text) and find inconsistencies
	type groupKey struct {
		modID      string
		sourceText string
	}

	type groupData struct {
		translations map[string]*translationInfoForTrivial // target_text -> info
		keys         []string
		statuses     map[string]bool // track unique statuses
	}

	groups := make(map[groupKey]*groupData)

	for _, pair := range pairs {
		key := groupKey{modID: pair.ModID, sourceText: pair.SourceText}
		if groups[key] == nil {
			groups[key] = &groupData{
				translations: make(map[string]*translationInfoForTrivial),
				keys:         []string{},
				statuses:     make(map[string]bool),
			}
		}
		if groups[key].translations[pair.TargetText] == nil {
			groups[key].translations[pair.TargetText] = &translationInfoForTrivial{count: 0, status: pair.Status}
		}
		groups[key].translations[pair.TargetText].count++
		groups[key].keys = append(groups[key].keys, pair.Key)
		groups[key].statuses[pair.Status] = true
	}

	// Find groups with multiple different translations
	var issues []ConsistencyIssue

	for key, data := range groups {
		if len(data.translations) <= 1 {
			continue
		}

		// Skip non-actionable source texts
		if isNonActionableSource(key.sourceText) {
			continue
		}

		// Skip if the difference is only official vs translated (intentional)
		if len(data.translations) == 2 && data.statuses["official"] && data.statuses["translated"] {
			// Check if one translation is official and another is translated
			hasOfficialOnly := false
			hasTranslatedOnly := false
			for _, info := range data.translations {
				if info.status == "official" {
					hasOfficialOnly = true
				} else if info.status == "translated" {
					hasTranslatedOnly = true
				}
			}
			if hasOfficialOnly && hasTranslatedOnly {
				continue // Skip official vs translated conflicts
			}
		}

		// Skip trivial differences (only whitespace/punctuation)
		if isTrivialDifference(data.translations) {
			continue
		}

		// Multiple different translations for the same source text
		issue := ConsistencyIssue{
			ModID:        key.modID,
			SourceText:   key.sourceText,
			Translations: make([]string, 0, len(data.translations)),
			Counts:       make(map[string]int),
			Keys:         data.keys,
		}

		// Sort translations by count (descending)
		type transCount struct {
			text  string
			count int
		}
		var sorted []transCount
		for t, info := range data.translations {
			sorted = append(sorted, transCount{t, info.count})
			issue.Translations = append(issue.Translations, t)
			issue.Counts[t] = info.count
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].count > sorted[j].count
		})

		// Suggest the most frequent translation
		issue.Suggested = sorted[0].text
		issue.Reason = "most frequent"

		// Check if there's an official translation
		officialTrans := a.findOfficialTranslation(ctx, key.modID, key.sourceText)
		if officialTrans != "" {
			issue.Suggested = officialTrans
			issue.Reason = "official translation"
		}

		// Check term dictionary
		termTrans := a.findTermTranslation(ctx, key.modID, key.sourceText)
		if termTrans != "" {
			issue.Suggested = termTrans
			issue.Reason = "term dictionary"
		}

		issues = append(issues, issue)
	}

	// Sort issues by number of affected keys (descending)
	sort.Slice(issues, func(i, j int) bool {
		return len(issues[i].Keys) > len(issues[j].Keys)
	})

	return issues, totalCount, nil
}

// translationInfoForTrivial is used for trivial difference check.
type translationInfoForTrivial struct {
	count  int
	status string
}

// isNonActionableSource checks if the source text is not worth tracking for consistency.
func isNonActionableSource(source string) bool {
	// Skip separator/comment lines
	if strings.HasPrefix(source, "===") || strings.HasPrefix(source, "---") {
		return true
	}
	if strings.Contains(source, "====") {
		return true
	}

	// Skip placeholder-only strings
	trimmed := strings.TrimSpace(source)
	if trimmed == "%s" || trimmed == "%d" || trimmed == "%1$s" || trimmed == "%2$s" {
		return true
	}

	// Skip very short strings (1-2 chars) that are likely punctuation/symbols
	if len(trimmed) <= 2 && !isJapaneseString(trimmed) {
		return true
	}

	return false
}

// isJapaneseString checks if a string contains Japanese characters.
func isJapaneseString(s string) bool {
	for _, r := range s {
		if r >= 0x3000 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

// isTrivialDifference checks if the translations only differ in whitespace or punctuation.
func isTrivialDifference(translations map[string]*translationInfoForTrivial) bool {
	if len(translations) != 2 {
		return false
	}

	var texts []string
	for t := range translations {
		texts = append(texts, t)
	}

	// Normalize both strings (remove spaces and certain punctuation)
	normalize := func(s string) string {
		result := strings.Builder{}
		for _, r := range s {
			if r != ' ' && r != '　' && r != '・' && r != ':' && r != '：' {
				result.WriteRune(r)
			}
		}
		return result.String()
	}

	return normalize(texts[0]) == normalize(texts[1])
}

// findOfficialTranslation looks for an official translation status.
func (a *Analyzer) findOfficialTranslation(ctx context.Context, modID, sourceText string) string {
	db := a.repo.GetDB()

	var result struct {
		TargetText string
	}

	err := db.Table("translations").
		Select("translations.target_text").
		Joins("JOIN translation_sources ON translation_sources.id = translations.source_id").
		Where("translation_sources.mod_id = ?", modID).
		Where("translation_sources.source_text = ?", sourceText).
		Where("translations.status = ?", "official").
		Limit(1).
		Scan(&result).Error

	if err != nil || result.TargetText == "" {
		return ""
	}

	return result.TargetText
}

// findTermTranslation looks for a matching term in the dictionary.
func (a *Analyzer) findTermTranslation(ctx context.Context, modID, sourceText string) string {
	// Get all terms and filter manually
	terms, err := a.repo.ListTerms(ctx, interfaces.TermFilter{})
	if err != nil || len(terms) == 0 {
		return ""
	}

	// Filter terms that match the source text
	var matchingTerms []*models.Term
	for _, term := range terms {
		if strings.EqualFold(term.SourceText, sourceText) {
			matchingTerms = append(matchingTerms, term)
		}
	}

	if len(matchingTerms) == 0 {
		return ""
	}

	terms = matchingTerms

	// Find the most specific matching term (mod > category > global)
	var bestMatch string
	bestPriority := -1

	for _, term := range terms {
		priority := 0
		switch {
		case term.Scope == "mod:"+modID:
			priority = 300
		case len(term.Scope) > 9 && term.Scope[:9] == "category:":
			priority = 200
		case term.Scope == "global":
			priority = 100
		}

		// Also consider the priority field
		priority += term.Priority

		if priority > bestPriority {
			bestPriority = priority
			bestMatch = term.TargetText
		}
	}

	return bestMatch
}
