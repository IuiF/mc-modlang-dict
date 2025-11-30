package jar

import (
	"regexp"
	"strings"
)

// MatchResult contains information about a pattern match.
type MatchResult struct {
	Path string            // Original file path
	Vars map[string]string // Extracted variables
}

// Matcher handles pattern matching for file paths.
type Matcher struct{}

// NewMatcher creates a new pattern matcher.
func NewMatcher() *Matcher {
	return &Matcher{}
}

// Match checks if a path matches a pattern and extracts variables.
// Pattern syntax:
//   - {name} - matches a single path segment and captures it
//   - * - matches a single path segment (wildcard)
//   - ** - matches zero or more path segments (recursive wildcard)
func (m *Matcher) Match(pattern, path string) (bool, map[string]string) {
	// Convert pattern to regex
	regex, varNames := m.patternToRegex(pattern)

	re, err := regexp.Compile("^" + regex + "$")
	if err != nil {
		return false, nil
	}

	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return false, nil
	}

	vars := make(map[string]string)
	for i, name := range varNames {
		if i+1 < len(matches) {
			vars[name] = matches[i+1]
		}
	}

	return true, vars
}

// FindFiles finds all files matching a pattern.
func (m *Matcher) FindFiles(pattern string, files []string) []MatchResult {
	var results []MatchResult

	for _, file := range files {
		if matched, vars := m.Match(pattern, file); matched {
			results = append(results, MatchResult{
				Path: file,
				Vars: vars,
			})
		}
	}

	return results
}

// ExpandPattern replaces variables in a pattern with their values.
func (m *Matcher) ExpandPattern(pattern string, vars map[string]string) string {
	result := pattern
	for name, value := range vars {
		result = strings.ReplaceAll(result, "{"+name+"}", value)
	}
	return result
}

// patternToRegex converts a pattern to a regex and returns variable names.
func (m *Matcher) patternToRegex(pattern string) (string, []string) {
	var regex strings.Builder
	var varNames []string

	parts := strings.Split(pattern, "/")

	for i, part := range parts {
		if i > 0 {
			regex.WriteString("/")
		}

		if part == "**" {
			// Match zero or more path segments
			regex.WriteString(".*")
		} else if part == "*" {
			// Match single path segment (non-greedy)
			regex.WriteString("[^/]+")
		} else {
			// Process each segment for {var} patterns
			segmentRegex, names := m.processSegment(part)
			regex.WriteString(segmentRegex)
			varNames = append(varNames, names...)
		}
	}

	return regex.String(), varNames
}

// processSegment converts a single path segment to regex.
func (m *Matcher) processSegment(segment string) (string, []string) {
	var regex strings.Builder
	var varNames []string

	i := 0
	for i < len(segment) {
		if segment[i] == '{' {
			// Find closing brace
			end := strings.Index(segment[i:], "}")
			if end == -1 {
				// No closing brace, treat as literal
				regex.WriteString(regexp.QuoteMeta(string(segment[i])))
				i++
				continue
			}

			varName := segment[i+1 : i+end]
			varNames = append(varNames, varName)
			regex.WriteString("([^/]+)")
			i += end + 1
		} else if segment[i] == '*' {
			// Wildcard within segment
			regex.WriteString("[^/]*")
			i++
		} else {
			// Literal character
			regex.WriteString(regexp.QuoteMeta(string(segment[i])))
			i++
		}
	}

	return regex.String(), varNames
}
