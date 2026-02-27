package domain

import (
	"testing"

	"github.com/annenpolka/cclog/internal/testutil"
)

func TestReplaceNewlinesWithSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no newlines", "hello world", "hello world"},
		{"LF", "hello\nworld", "hello world"},
		{"CR", "hello\rworld", "hello world"},
		{"CRLF", "hello\r\nworld", "hello world"},
		{"multiple LF", "a\nb\nc", "a b c"},
		{"mixed newlines", "a\nb\r\nc\rd", "a b c d"},
		{"empty string", "", ""},
		{"only newlines", "\n\n\n", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.Diff(t, tt.expected, ReplaceNewlinesWithSpaces(tt.input))
		})
	}
}
