package analyzer

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"github.com/iuif/minecraft-mod-dictionary/pkg/interfaces"
)

// phraseData holds information about a discovered phrase.
type phraseData struct {
	phrase       string
	translations map[string]int   // translation -> count
	examples     []PhraseExample
}

// minePhrases discovers phrase patterns through N-gram mining.
func (a *Analyzer) minePhrases(ctx context.Context, modID string, minCount int) ([]PhraseIssue, error) {
	// Get all translation pairs
	pairs, err := a.getAllTranslationPairs(ctx, modID)
	if err != nil {
		return nil, err
	}

	// Extract N-grams from source texts and build phrase-to-translations mapping
	phraseMap := make(map[string]*phraseData)

	for _, pair := range pairs {
		// Extract N-grams (1 to 4 words)
		ngrams := extractNgrams(pair.SourceText, 1, 4)

		for _, ngram := range ngrams {
			if phraseMap[ngram] == nil {
				phraseMap[ngram] = &phraseData{
					phrase:       ngram,
					translations: make(map[string]int),
					examples:     []PhraseExample{},
				}
			}

			// Try to find the corresponding translation for this N-gram
			trans := findCorrespondingTranslation(pair.SourceText, pair.TargetText, ngram)
			if trans != "" {
				phraseMap[ngram].translations[trans]++

				// Store example (limit to avoid memory issues)
				if len(phraseMap[ngram].examples) < 10 {
					phraseMap[ngram].examples = append(phraseMap[ngram].examples, PhraseExample{
						Key:        pair.Key,
						SourceText: pair.SourceText,
						TargetText: pair.TargetText,
					})
				}
			}
		}
	}

	// Get terms for dictionary lookup
	terms, _ := a.repo.ListTerms(ctx, interfaces.TermFilter{})
	termMap := make(map[string]string)
	for _, t := range terms {
		termMap[strings.ToLower(t.SourceText)] = t.TargetText
	}

	// Convert to PhraseIssue list and filter by minCount
	var issues []PhraseIssue

	for phrase, data := range phraseMap {
		// Calculate total occurrences
		totalOccurrences := 0
		for _, count := range data.translations {
			totalOccurrences += count
		}

		// Skip if below minimum count
		if totalOccurrences < minCount {
			continue
		}

		// Skip single-word phrases that are too common or too short
		words := strings.Fields(phrase)
		if len(words) == 1 {
			// Skip very short words
			if len(phrase) < 3 {
				continue
			}
			// Skip common words
			if isCommonWord(phrase) {
				continue
			}
		}

		// Calculate consistency score
		maxCount := 0
		var suggested string
		for trans, count := range data.translations {
			if count > maxCount {
				maxCount = count
				suggested = trans
			}
		}

		consistencyScore := float64(maxCount) / float64(totalOccurrences)

		// Skip if perfectly consistent (score = 1.0)
		if consistencyScore >= 1.0 {
			continue
		}

		// Check if in dictionary
		inDict := false
		if dictTrans, ok := termMap[strings.ToLower(phrase)]; ok {
			inDict = true
			// If dictionary has different suggestion, prefer it
			if dictTrans != suggested {
				suggested = dictTrans
			}
		}

		issue := PhraseIssue{
			Phrase:           phrase,
			Occurrences:      totalOccurrences,
			Translations:     data.translations,
			ConsistencyScore: consistencyScore,
			InDictionary:     inDict,
			Suggested:        suggested,
			Examples:         data.examples,
		}

		issues = append(issues, issue)
	}

	// Sort by: 1) inconsistency (lower score first), 2) occurrences (higher first)
	sort.Slice(issues, func(i, j int) bool {
		// First, prioritize inconsistent phrases
		if issues[i].ConsistencyScore != issues[j].ConsistencyScore {
			return issues[i].ConsistencyScore < issues[j].ConsistencyScore
		}
		// Then by occurrences
		return issues[i].Occurrences > issues[j].Occurrences
	})

	return issues, nil
}

// getAllTranslationPairs retrieves all translation pairs from the database.
func (a *Analyzer) getAllTranslationPairs(ctx context.Context, modID string) ([]SourceTranslationPair, error) {
	db := a.repo.GetDB()

	var pairs []SourceTranslationPair

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
		return nil, err
	}

	return pairs, nil
}

// extractNgrams extracts word N-grams from text.
func extractNgrams(text string, minN, maxN int) []string {
	// Tokenize: split by spaces and punctuation
	words := tokenize(text)
	if len(words) == 0 {
		return nil
	}

	var ngrams []string
	seen := make(map[string]bool)

	for n := minN; n <= maxN; n++ {
		if n > len(words) {
			break
		}

		for i := 0; i <= len(words)-n; i++ {
			ngram := strings.Join(words[i:i+n], " ")
			if !seen[ngram] {
				seen[ngram] = true
				ngrams = append(ngrams, ngram)
			}
		}
	}

	return ngrams
}

// tokenize splits text into words, handling punctuation.
func tokenize(text string) []string {
	var words []string
	var current strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' || r == '-' {
			current.WriteRune(r)
		} else if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// findCorrespondingTranslation attempts to find the Japanese translation for an English phrase.
// This is a heuristic approach based on position and pattern matching.
func findCorrespondingTranslation(sourceText, targetText, phrase string) string {
	// Simple approach: if the phrase is the entire source, return the entire target
	if strings.EqualFold(sourceText, phrase) {
		return targetText
	}

	// Check if the phrase appears in the source
	phraseIdx := strings.Index(strings.ToLower(sourceText), strings.ToLower(phrase))
	if phraseIdx == -1 {
		return ""
	}

	// Calculate relative position in source
	sourceLen := len(sourceText)
	if sourceLen == 0 {
		return ""
	}

	relativeStart := float64(phraseIdx) / float64(sourceLen)
	relativeEnd := float64(phraseIdx+len(phrase)) / float64(sourceLen)

	// Apply the same relative position to target
	targetLen := len(targetText)
	if targetLen == 0 {
		return ""
	}

	// For Japanese text, we need to be more careful with character boundaries
	targetRunes := []rune(targetText)
	targetRuneLen := len(targetRunes)

	startIdx := int(relativeStart * float64(targetRuneLen))
	endIdx := int(relativeEnd * float64(targetRuneLen))

	// Clamp to valid range
	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > targetRuneLen {
		endIdx = targetRuneLen
	}
	if startIdx >= endIdx {
		return ""
	}

	// Try to find word boundaries in Japanese
	// Extend to nearest particle or punctuation
	for startIdx > 0 && !isJapaneseBreak(targetRunes[startIdx-1]) {
		startIdx--
	}
	for endIdx < targetRuneLen && !isJapaneseBreak(targetRunes[endIdx]) {
		endIdx++
	}

	result := string(targetRunes[startIdx:endIdx])
	result = strings.TrimSpace(result)
	result = strings.Trim(result, "のをがはにでとからまでより")

	if len(result) == 0 {
		return ""
	}

	return result
}

// isJapaneseBreak checks if a rune is a natural break point in Japanese text.
func isJapaneseBreak(r rune) bool {
	// Particles and punctuation
	breaks := "のをがはにでとからまでより・、。「」『』（）【】"
	return strings.ContainsRune(breaks, r) || unicode.IsSpace(r)
}

// isCommonWord checks if a word is too common to be interesting.
func isCommonWord(word string) bool {
	common := map[string]bool{
		// English articles and prepositions
		"the": true, "a": true, "an": true, "of": true, "to": true,
		"and": true, "or": true, "in": true, "on": true, "at": true,
		"for": true, "with": true, "by": true, "from": true, "as": true,
		"is": true, "are": true, "was": true, "were": true, "be": true,
		"has": true, "have": true, "had": true, "do": true, "does": true,
		"this": true, "that": true, "these": true, "those": true,
		"it": true, "its": true, "i": true, "you": true, "we": true,
		"they": true, "he": true, "she": true, "my": true, "your": true,
		"will": true, "would": true, "can": true, "could": true,
		"should": true, "may": true, "might": true, "must": true,
		"all": true, "any": true, "some": true, "no": true, "not": true,
		// Minecraft colors (often translated differently by context)
		"white": true, "orange": true, "magenta": true, "light": true,
		"yellow": true, "lime": true, "pink": true, "gray": true,
		"cyan": true, "purple": true, "blue": true, "brown": true,
		"green": true, "red": true, "black": true, "grey": true,
		// Common Minecraft terms that vary by context
		"block": true, "item": true, "tile": true, "entity": true,
		"slab": true, "stairs": true, "wall": true, "fence": true,
		"button": true, "plate": true, "door": true, "gate": true,
		"small": true, "large": true, "big": true, "tiny": true,
	}
	return common[strings.ToLower(word)]
}
