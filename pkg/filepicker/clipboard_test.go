package filepicker

import (
	"testing"

	"github.com/atotto/clipboard"

	"github.com/annenpolka/cclog/internal/testutil"
)

func skipWithoutClipboard(t *testing.T) {
	t.Helper()
	if err := clipboard.WriteAll("test"); err != nil {
		t.Skip("clipboard not available in this environment")
	}
}

func TestCopySessionID(t *testing.T) {
	skipWithoutClipboard(t)

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid jsonl file",
			filePath: "../../testdata/sample.jsonl",
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "non-existent.jsonl",
			wantErr:  false, // extractSessionID doesn't check file existence, only filename format
		},
		{
			name:     "empty file path",
			filePath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the copySessionID command
			cmd := copySessionID(tt.filePath)
			msg := cmd()

			// Check if the result is of the expected type
			result, ok := msg.(copySessionIDMsg)
			testutil.True(t, ok)
			// Check if error expectation matches
			testutil.Diff(t, tt.wantErr, result.error != nil)
		})
	}
}

func TestCopySessionIDIntegration(t *testing.T) {
	skipWithoutClipboard(t)

	filePath := "../../testdata/sample.jsonl"

	// Execute the copySessionID command
	cmd := copySessionID(filePath)
	msg := cmd()

	// Check if the result is of the expected type
	result, ok := msg.(copySessionIDMsg)
	testutil.True(t, ok)
	// For a valid file, we should get a successful result
	testutil.Diff(t, false, result.error != nil)
	testutil.True(t, result.success)
}

func TestCopySessionIDErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "invalid file extension",
			filePath: "test.txt",
			wantErr:  true,
		},
		{
			name:     "file without extension",
			filePath: "test",
			wantErr:  true,
		},
		{
			name:     "only extension",
			filePath: ".jsonl",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the copySessionID command
			cmd := copySessionID(tt.filePath)
			msg := cmd()

			// Check if the result is of the expected type
			result, ok := msg.(copySessionIDMsg)
			testutil.True(t, ok)
			// Check if error expectation matches
			testutil.Diff(t, tt.wantErr, result.error != nil)
		})
	}
}
