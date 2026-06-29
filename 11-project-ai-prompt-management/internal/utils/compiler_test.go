package utils

import (
	"testing"
)

func TestCompilePrompt(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		vars          map[string]string
		expected      string
		expectedToken int
	}{
		{
			name:     "Simple variable replacement",
			template: "Hello {{name}}!",
			vars:     map[string]string{"name": "Timur"},
			expected: "Hello Timur!",
			// "Hello Timur!" has 2 words. 2 * 1.33 = 2.66 => int 2
			expectedToken: 2,
		},
		{
			name:          "Variable with spaces inside braces",
			template:      "Welcome to {{  platform  }} dev server.",
			vars:          map[string]string{"platform": "Gemini"},
			expected:      "Welcome to Gemini dev server.",
			expectedToken: 6, // 5 words * 1.33 = 6.65 => 6
		},
		{
			name:          "Missing variable replaced with empty string",
			template:      "User {{username}} logged in from {{ip}}.",
			vars:          map[string]string{"username": "antigravity"},
			expected:      "User antigravity logged in from .",
			expectedToken: 7, // 6 words * 1.33 = 7.98 => 7
		},
		{
			name:          "Multiple variables replacement",
			template:      "Query: {{query}} in category {{category}} by user {{user}}.",
			vars:          map[string]string{"query": "golang", "category": "programming", "user": "timur"},
			expected:      "Query: golang in category programming by user timur.",
			expectedToken: 10, // 8 words * 1.33 = 10.64 => 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotToken := CompilePrompt(tt.template, tt.vars)
			if gotText != tt.expected {
				t.Errorf("CompilePrompt() gotText = %q, expected %q", gotText, tt.expected)
			}
			if gotToken != tt.expectedToken {
				t.Errorf("CompilePrompt() gotToken = %d, expected %d", gotToken, tt.expectedToken)
			}
		})
	}
}
