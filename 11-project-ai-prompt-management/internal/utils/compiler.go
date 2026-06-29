package utils

import (
	"regexp"
	"strings"
)

var variableRegex = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_]+)\s*\}\}`)

// CompilePrompt parses a prompt template string containing double curly braces (e.g. {{name}})
// and replaces them with values provided in the vars map.
// It also returns an estimated token count (using industrial standard words * 1.33 calculation).
func CompilePrompt(templateText string, vars map[string]string) (string, int) {
	compiled := variableRegex.ReplaceAllStringFunc(templateText, func(match string) string {
		// Extract variable name from {{var}}
		subMatches := variableRegex.FindStringSubmatch(match)
		if len(subMatches) < 2 {
			return ""
		}
		varName := subMatches[1]
		val, exists := vars[varName]
		if !exists {
			return "" // Fallback to empty string if variable is not provided
		}
		return val
	})

	// Estimate token length based on word count
	words := strings.Fields(compiled)
	wordCount := len(words)
	// Typical estimation: 1 word ~ 1.33 tokens
	tokenEstimate := int(float64(wordCount) * 1.33)
	if tokenEstimate == 0 && len(compiled) > 0 {
		tokenEstimate = 1
	}

	return compiled, tokenEstimate
}
