package analyzer

import (
	"context"
	"sort"
	"time"

	"github.com/iuif/minecraft-mod-dictionary/internal/database"
)

// Analyzer provides translation consistency analysis.
type Analyzer struct {
	repo *database.Repository
}

// New creates a new Analyzer.
func New(repo *database.Repository) *Analyzer {
	return &Analyzer{repo: repo}
}

// AnalyzeAll runs all analysis types.
func (a *Analyzer) AnalyzeAll(ctx context.Context, opts AnalysisOptions) (*AnalysisResult, error) {
	result := &AnalysisResult{
		AnalysisDate: time.Now(),
		TargetMod:    opts.ModID,
	}

	// Consistency analysis
	consistencyResult, err := a.AnalyzeConsistency(ctx, opts)
	if err != nil {
		return nil, err
	}
	result.Consistency = consistencyResult.Consistency
	result.Summary.ConsistencyIssues = len(result.Consistency)

	// Phrase analysis (with stricter filtering for 'all' command)
	phraseOpts := opts
	if phraseOpts.MinCount < 10 {
		phraseOpts.MinCount = 10 // Higher threshold for all command
	}
	phrasesResult, err := a.AnalyzePhrases(ctx, phraseOpts)
	if err != nil {
		return nil, err
	}
	// Only include phrases that are significantly inconsistent (score < 0.5)
	// and have high occurrences - limit to top 500 most impactful
	var candidatePhrases []PhraseIssue
	for _, p := range phrasesResult.Phrases {
		if p.ConsistencyScore < 0.5 && p.Occurrences >= 15 {
			candidatePhrases = append(candidatePhrases, p)
		}
	}
	// Sort by impact: lower consistency + higher occurrences = higher priority
	sort.Slice(candidatePhrases, func(i, j int) bool {
		// Prioritize: more inconsistent first, then more occurrences
		impactI := float64(candidatePhrases[i].Occurrences) * (1.0 - candidatePhrases[i].ConsistencyScore)
		impactJ := float64(candidatePhrases[j].Occurrences) * (1.0 - candidatePhrases[j].ConsistencyScore)
		return impactI > impactJ
	})
	// Limit to top 500
	maxPhrases := 500
	if len(candidatePhrases) > maxPhrases {
		candidatePhrases = candidatePhrases[:maxPhrases]
	}
	result.Phrases = candidatePhrases
	result.Summary.DiscoveredPhrases = len(result.Phrases)

	// Count inconsistent phrases
	for _, p := range result.Phrases {
		if p.ConsistencyScore < 1.0 {
			result.Summary.InconsistentPhrases++
		}
	}

	// Terms analysis
	termsResult, err := a.AnalyzeTerms(ctx, opts)
	if err != nil {
		return nil, err
	}
	result.TermViolations = termsResult.TermViolations
	result.Summary.TermViolations = len(result.TermViolations)

	// Get total translation count
	result.Summary.TotalTranslations = consistencyResult.Summary.TotalTranslations

	return result, nil
}

// AnalyzeConsistency checks for same source text with different translations.
func (a *Analyzer) AnalyzeConsistency(ctx context.Context, opts AnalysisOptions) (*AnalysisResult, error) {
	result := &AnalysisResult{
		AnalysisDate: time.Now(),
		TargetMod:    opts.ModID,
	}

	issues, totalCount, err := a.findConsistencyIssues(ctx, opts.ModID)
	if err != nil {
		return nil, err
	}

	result.Consistency = issues
	result.Summary.ConsistencyIssues = len(issues)
	result.Summary.TotalTranslations = totalCount

	return result, nil
}

// AnalyzePhrases discovers phrase patterns through N-gram mining.
func (a *Analyzer) AnalyzePhrases(ctx context.Context, opts AnalysisOptions) (*AnalysisResult, error) {
	result := &AnalysisResult{
		AnalysisDate: time.Now(),
		TargetMod:    opts.ModID,
	}

	phrases, err := a.minePhrases(ctx, opts.ModID, opts.MinCount)
	if err != nil {
		return nil, err
	}

	result.Phrases = phrases
	result.Summary.DiscoveredPhrases = len(phrases)

	for _, p := range phrases {
		if p.ConsistencyScore < 1.0 {
			result.Summary.InconsistentPhrases++
		}
	}

	return result, nil
}

// AnalyzeTerms checks translations against term dictionary.
func (a *Analyzer) AnalyzeTerms(ctx context.Context, opts AnalysisOptions) (*AnalysisResult, error) {
	result := &AnalysisResult{
		AnalysisDate: time.Now(),
		TargetMod:    opts.ModID,
	}

	violations, err := a.checkTermCompliance(ctx, opts.ModID)
	if err != nil {
		return nil, err
	}

	result.TermViolations = violations
	result.Summary.TermViolations = len(violations)

	return result, nil
}
