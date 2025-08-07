package filepicker

import (
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"testing"

	"github.com/annenpolka/cclog/internal/testutil"
)

func TestPreviewModel_SetContent(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectedContent string
	}{
		{
			name:            "Empty content",
			content:         "",
			expectedContent: "",
		},
		{
			name:            "Simple markdown content",
			content:         "# Title\n\nThis is a test.",
			expectedContent: "# Title\n\nThis is a test.",
		},
		{
			name:            "Multi-line content",
			content:         "Line 1\nLine 2\nLine 3",
			expectedContent: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preview := NewPreviewModel()
			_ = preview.SetContent(tt.content)
			testutil.Diff(t, tt.expectedContent, preview.GetContent())
		})
	}
}

func TestPreviewModel_SetVisible(t *testing.T) {
	tests := []struct {
		name    string
		visible bool
	}{
		{
			name:    "Set visible to true",
			visible: true,
		},
		{
			name:    "Set visible to false",
			visible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preview := NewPreviewModel()
			preview.SetVisible(tt.visible)
			testutil.Diff(t, tt.visible, preview.IsVisible())
		})
	}
}

func TestPreviewModel_SetSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{
			name:   "Set size to 80x24",
			width:  80,
			height: 24,
		},
		{
			name:   "Set size to 120x40",
			width:  120,
			height: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preview := NewPreviewModel()
			preview.SetSize(tt.width, tt.height)

			width, height := preview.GetSize()
			testutil.Diff(t, tt.width, width)
			testutil.Diff(t, tt.height, height)
		})
	}
}

func TestGeneratePreview(t *testing.T) {
	tests := []struct {
		name          string
		jsonlPath     string
		shouldError   bool
		expectedEmpty bool
	}{
		{
			name:          "Valid JSONL file",
			jsonlPath:     "../../testdata/sample.jsonl",
			shouldError:   false,
			expectedEmpty: false,
		},
		{
			name:          "Non-existent file",
			jsonlPath:     "non-existent-file.jsonl",
			shouldError:   true,
			expectedEmpty: true,
		},
		{
			name:          "Empty path",
			jsonlPath:     "",
			shouldError:   false,
			expectedEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GeneratePreview(tt.jsonlPath, true)
			testutil.Diff(t, tt.shouldError, err != nil)
			if tt.expectedEmpty {
				testutil.Diff(t, "", content)
			} else if !tt.shouldError {
				testutil.Diff(t, false, content == "")
			}
		})
	}
}

func TestPreviewModel_DefaultState(t *testing.T) {
	preview := NewPreviewModel()

	testutil.True(t, preview.IsVisible())

	testutil.Diff(t, "", preview.GetContent())

	width, height := preview.GetSize()
	testutil.Diff(t, 0, width)
	testutil.Diff(t, 0, height)
}

func TestPreviewModel_Cleanup(t *testing.T) {
	preview := NewPreviewModel()

	// Set some content to create temp file
	_ = preview.SetContent("# Test Content\n\nThis is a test.")

	// Check that temp file was created
	testutil.Diff(t, false, preview.tempFile == "")

	// Check temp file exists
	_, err := os.Stat(preview.tempFile)
	testutil.Diff(t, false, os.IsNotExist(err))

	// Cleanup should remove temp file
	preview.Cleanup()
	// Check temp file is removed
	testutil.Diff(t, "", preview.tempFile)
}

func TestPreviewModel_KeyBindings_GoToTop(t *testing.T) {
	preview := NewPreviewModel()

	// Set some content and scroll position
	cmd := preview.SetContent("# Test Content\n\nLine 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6")
	if cmd != nil {
		// Execute the command to load content
		_ = cmd()
	}
	preview.SetSize(80, 10)

	// Simulate scrolling down first (so we have somewhere to scroll back to)
	preview.markdownBubble.Viewport.ScrollDown(5)

	// Simulate 'g' key press
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	preview.Update(keyMsg)

	// Check that we're now at the top
	testutil.Diff(t, 0, preview.markdownBubble.Viewport.YOffset)
}

func TestPreviewModel_KeyBindings_GoToBottom(t *testing.T) {
	preview := NewPreviewModel()

	// Set some content that will be longer than the viewport
	longContent := "# Test Content\n\n"
	for i := 0; i < 20; i++ {
		longContent += "Line " + string(rune('A'+i)) + "\n"
	}
	cmd := preview.SetContent(longContent)
	if cmd != nil {
		// Execute the command to load content
		_ = cmd()
	}
	preview.SetSize(80, 10)

	// Initially should be at top
	testutil.Diff(t, 0, preview.markdownBubble.Viewport.YOffset)

	// Simulate 'G' key press (shift+g)
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	preview.Update(keyMsg)

	// Check that we're now at the bottom
	// The bottom position should be greater than 0 for content longer than viewport
	finalOffset := preview.markdownBubble.Viewport.YOffset
	totalLines := preview.markdownBubble.Viewport.TotalLineCount()
	height := preview.markdownBubble.Viewport.Height

	// For content that's longer than viewport, we should have scrolled down
	if totalLines > height {
		testutil.Diff(t, false, finalOffset == 0)
	}
}
