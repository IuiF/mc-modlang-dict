package analyzer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// FormatJSON formats the result as JSON.
func FormatJSON(result *AnalysisResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

// FormatCSV formats the result as CSV.
func FormatCSV(result *AnalysisResult) ([]byte, error) {
	var buf strings.Builder
	w := csv.NewWriter(&buf)

	// Header
	header := []string{"type", "mod_id", "phrase", "translations", "suggested", "occurrences", "consistency_score", "in_dictionary"}
	if err := w.Write(header); err != nil {
		return nil, err
	}

	// Consistency issues
	for _, issue := range result.Consistency {
		translations := strings.Join(issue.Translations, "|")
		record := []string{
			"consistency",
			issue.ModID,
			issue.SourceText,
			translations,
			issue.Suggested,
			fmt.Sprintf("%d", len(issue.Keys)),
			"",
			"",
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	// Discovered phrases
	for _, phrase := range result.Phrases {
		var translations []string
		for t := range phrase.Translations {
			translations = append(translations, t)
		}
		sort.Strings(translations)

		inDict := "false"
		if phrase.InDictionary {
			inDict = "true"
		}

		record := []string{
			"discovered",
			result.TargetMod,
			phrase.Phrase,
			strings.Join(translations, "|"),
			phrase.Suggested,
			fmt.Sprintf("%d", phrase.Occurrences),
			fmt.Sprintf("%.2f", phrase.ConsistencyScore),
			inDict,
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	// Term violations
	for _, violation := range result.TermViolations {
		record := []string{
			"term_violation",
			result.TargetMod,
			violation.TermSource,
			"",
			violation.TermTarget,
			fmt.Sprintf("%d", violation.ViolationCount),
			"0.0",
			"true",
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return []byte(buf.String()), nil
}

// FormatSummary formats the result as human-readable summary.
func FormatSummary(result *AnalysisResult) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString("=== Translation Consistency Analysis ===\n")
	if result.TargetMod != "" {
		buf.WriteString(fmt.Sprintf("Target: %s\n", result.TargetMod))
	} else {
		buf.WriteString("Target: All mods\n")
	}
	buf.WriteString(fmt.Sprintf("Date: %s\n\n", result.AnalysisDate.Format("2006-01-02 15:04:05")))

	// Consistency issues
	if len(result.Consistency) > 0 {
		buf.WriteString(fmt.Sprintf("--- Consistency Issues (%d) ---\n", len(result.Consistency)))
		for _, issue := range result.Consistency {
			buf.WriteString(fmt.Sprintf("\n[ISSUE] \"%s\"\n", issue.SourceText))
			buf.WriteString(fmt.Sprintf("  Mod: %s\n", issue.ModID))
			buf.WriteString("  Translations:\n")
			for _, t := range issue.Translations {
				count := issue.Counts[t]
				buf.WriteString(fmt.Sprintf("    - \"%s\" (%d)\n", t, count))
			}
			buf.WriteString(fmt.Sprintf("  Suggested: \"%s\"\n", issue.Suggested))
			if issue.Reason != "" {
				buf.WriteString(fmt.Sprintf("  Reason: %s\n", issue.Reason))
			}
		}
		buf.WriteString("\n")
	}

	// Discovered phrases
	if len(result.Phrases) > 0 {
		buf.WriteString(fmt.Sprintf("--- Discovered Phrases (%d) ---\n", len(result.Phrases)))

		// Sort by consistency score (lower = more inconsistent = higher priority)
		sortedPhrases := make([]PhraseIssue, len(result.Phrases))
		copy(sortedPhrases, result.Phrases)
		sort.Slice(sortedPhrases, func(i, j int) bool {
			// First by consistency (lower first), then by occurrences (higher first)
			if sortedPhrases[i].ConsistencyScore != sortedPhrases[j].ConsistencyScore {
				return sortedPhrases[i].ConsistencyScore < sortedPhrases[j].ConsistencyScore
			}
			return sortedPhrases[i].Occurrences > sortedPhrases[j].Occurrences
		})

		for _, phrase := range sortedPhrases {
			dictMark := ""
			if !phrase.InDictionary {
				dictMark = " [NEW]"
			}

			buf.WriteString(fmt.Sprintf("\n[%.0f%%] \"%s\"%s\n", phrase.ConsistencyScore*100, phrase.Phrase, dictMark))
			buf.WriteString(fmt.Sprintf("  Occurrences: %d\n", phrase.Occurrences))
			buf.WriteString("  Translations:\n")

			// Sort translations by count
			type transCount struct {
				text  string
				count int
			}
			var trans []transCount
			for t, c := range phrase.Translations {
				trans = append(trans, transCount{t, c})
			}
			sort.Slice(trans, func(i, j int) bool {
				return trans[i].count > trans[j].count
			})

			for _, t := range trans {
				buf.WriteString(fmt.Sprintf("    - \"%s\" (%d)\n", t.text, t.count))
			}
			buf.WriteString(fmt.Sprintf("  Suggested: \"%s\"\n", phrase.Suggested))

			// Show examples if available
			if len(phrase.Examples) > 0 {
				buf.WriteString("  Examples:\n")
				maxExamples := 3
				if len(phrase.Examples) < maxExamples {
					maxExamples = len(phrase.Examples)
				}
				for i := 0; i < maxExamples; i++ {
					ex := phrase.Examples[i]
					buf.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\"\n", ex.SourceText, ex.TargetText))
				}
			}
		}
		buf.WriteString("\n")
	}

	// Term violations
	if len(result.TermViolations) > 0 {
		buf.WriteString(fmt.Sprintf("--- Term Violations (%d) ---\n", len(result.TermViolations)))
		for _, v := range result.TermViolations {
			buf.WriteString(fmt.Sprintf("\n[VIOLATION] \"%s\" -> \"%s\"\n", v.TermSource, v.TermTarget))
			buf.WriteString(fmt.Sprintf("  Scope: %s\n", v.TermScope))
			buf.WriteString(fmt.Sprintf("  Violations: %d keys\n", v.ViolationCount))
			if len(v.Violations) > 0 {
				maxShow := 3
				if len(v.Violations) < maxShow {
					maxShow = len(v.Violations)
				}
				for i := 0; i < maxShow; i++ {
					d := v.Violations[i]
					buf.WriteString(fmt.Sprintf("    - \"%s\" -> \"%s\" (expected: \"%s\")\n", d.SourceText, d.TargetText, d.Expected))
				}
				if len(v.Violations) > maxShow {
					buf.WriteString(fmt.Sprintf("    ... and %d more\n", len(v.Violations)-maxShow))
				}
			}
		}
		buf.WriteString("\n")
	}

	// Summary
	buf.WriteString("=== Summary ===\n")
	buf.WriteString(fmt.Sprintf("Total translations: %d\n", result.Summary.TotalTranslations))
	buf.WriteString(fmt.Sprintf("Consistency issues: %d\n", result.Summary.ConsistencyIssues))
	buf.WriteString(fmt.Sprintf("Discovered phrases: %d\n", result.Summary.DiscoveredPhrases))
	buf.WriteString(fmt.Sprintf("Inconsistent phrases: %d\n", result.Summary.InconsistentPhrases))
	buf.WriteString(fmt.Sprintf("Term violations: %d\n", result.Summary.TermViolations))

	return []byte(buf.String()), nil
}
