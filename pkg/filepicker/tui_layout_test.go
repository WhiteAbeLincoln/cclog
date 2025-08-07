package filepicker

import (
	"testing"

	"github.com/annenpolka/cclog/internal/testutil"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdatePreviewSize(t *testing.T) {
	model := NewModel("/tmp", false)

	tests := []struct {
		name                  string
		terminalWidth         int
		terminalHeight        int
		expectedPreviewWidth  int
		expectedPreviewHeight int
	}{
		{
			name:                  "Standard terminal size",
			terminalWidth:         80,
			terminalHeight:        40,
			expectedPreviewWidth:  80, // Use full width
			expectedPreviewHeight: 27, // (40 - 6) * 0.8 = 27.2 -> 27
		},
		{
			name:                  "Large terminal",
			terminalWidth:         120,
			terminalHeight:        60,
			expectedPreviewWidth:  120, // Use full width
			expectedPreviewHeight: 43,  // (60 - 6) * 0.8 = 43.2 -> 43
		},
		{
			name:                  "Small terminal",
			terminalWidth:         40,
			terminalHeight:        20,
			expectedPreviewWidth:  40, // Use full width
			expectedPreviewHeight: 10, // Adaptive split gives more space to list
		},
		{
			name:                  "Very small terminal",
			terminalWidth:         10,
			terminalHeight:        8,
			expectedPreviewWidth:  10, // Use full width
			expectedPreviewHeight: 0,  // Prioritize list on very small screens
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.terminalWidth = tt.terminalWidth
			model.terminalHeight = tt.terminalHeight

			model.updatePreviewSize()

			width, height := model.preview.GetSize()

			testutil.Diff(t, tt.expectedPreviewWidth, width)
			testutil.Diff(t, tt.expectedPreviewHeight, height)
		})
	}
}

func TestModelDynamicLayoutAdjustment(t *testing.T) {
	model := NewModel("/tmp", false)

	tests := []struct {
		name                  string
		terminalWidth         int
		terminalHeight        int
		splitRatio            float64
		expectedListHeight    int
		expectedPreviewHeight int
	}{
		{
			name:                  "50/50 split with medium terminal",
			terminalWidth:         80,
			terminalHeight:        40,
			splitRatio:            0.5,
			expectedListHeight:    17, // (40 - 6) / 2
			expectedPreviewHeight: 17,
		},
		{
			name:                  "30/70 split favoring preview",
			terminalWidth:         80,
			terminalHeight:        60,
			splitRatio:            0.7,
			expectedListHeight:    17, // 60 - 6 - 37 = 17
			expectedPreviewHeight: 37, // (60 - 6) * 0.7 = 37.8 -> 37
		},
		{
			name:                  "70/30 split favoring list",
			terminalWidth:         80,
			terminalHeight:        40,
			splitRatio:            0.3,
			expectedListHeight:    24, // (40 - 6) * 0.7
			expectedPreviewHeight: 10, // (40 - 6) * 0.3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.terminalWidth = tt.terminalWidth
			model.terminalHeight = tt.terminalHeight

			model.updateDynamicLayout(tt.splitRatio)

			listHeight := model.getListHeight()
			_, previewHeight := model.preview.GetSize()

			testutil.Diff(t, tt.expectedListHeight, listHeight)
			testutil.Diff(t, tt.expectedPreviewHeight, previewHeight)
		})
	}
}

func TestModelWindowSizeMessage(t *testing.T) {
	model := NewModel("/tmp", false)

	// Test window size update
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)

	m := updatedModel.(Model)

	testutil.Diff(t, 100, m.terminalWidth)
	testutil.Diff(t, 50, m.terminalHeight)

	// Check if preview size was updated
	width, height := m.preview.GetSize()
	expectedWidth := 100 // Use full width
	expectedHeight := 35 // (50 - 6) * 0.8 = 35.2 -> 35

	testutil.Diff(t, expectedWidth, width)
	testutil.Diff(t, expectedHeight, height)
}

func TestModelResponsiveLayoutThresholds(t *testing.T) {
	model := NewModel("/tmp", false)

	tests := []struct {
		name                  string
		terminalWidth         int
		terminalHeight        int
		expectedCompact       bool
		expectedMaxTitleChars int
	}{
		{
			name:                  "Large terminal - full layout",
			terminalWidth:         120,
			terminalHeight:        60,
			expectedCompact:       false,
			expectedMaxTitleChars: 108, // 120 - 17 - 3 - 2 + (120-80)/4 = 108
		},
		{
			name:                  "Medium terminal - full layout",
			terminalWidth:         80,
			terminalHeight:        40,
			expectedCompact:       false,
			expectedMaxTitleChars: 58, // 80 - 17 - 3 - 2 = 58
		},
		{
			name:                  "Small terminal - compact layout",
			terminalWidth:         50,
			terminalHeight:        25,
			expectedCompact:       true,
			expectedMaxTitleChars: 28, // 50 - 17 - 3 - 2 = 28
		},
		{
			name:                  "Very small terminal - compact layout",
			terminalWidth:         30,
			terminalHeight:        15,
			expectedCompact:       true,
			expectedMaxTitleChars: 20, // Minimum threshold
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.terminalWidth = tt.terminalWidth
			model.terminalHeight = tt.terminalHeight

			model.updateDisplaySettings()

			testutil.Diff(t, tt.expectedCompact, model.useCompactLayout)
			testutil.Diff(t, tt.expectedMaxTitleChars, model.maxTitleChars)
		})
	}
}

func TestModelMinimumSizeConstraints(t *testing.T) {
	model := NewModel("/tmp", false)

	// Test with extremely small terminal
	model.terminalWidth = 5
	model.terminalHeight = 5

	model.updatePreviewSize()

	width, height := model.preview.GetSize()

	// Should handle negative sizes gracefully
	testutil.True(t, width >= 0)
	testutil.True(t, height >= 0)
}
