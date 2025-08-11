package cli

import (
	"fmt"

	"github.com/annenpolka/cclog/internal/usecase"
	"github.com/annenpolka/cclog/pkg/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// RunTUI starts the TUI file picker and returns the selected file
func RunTUI(config Config) (string, error) {
	// Inject renderer so TUI uses the same pipeline as CLI/usecase
	filepicker.SetRenderer(func(path string, includeAll bool) (string, error) {
		opts := usecase.Options{IncludeAll: includeAll, ShowUUID: false, ShowTitle: false}
		return usecase.GenerateMarkdownFromPath(path, opts)
	})

	// Create and run the TUI model
	model := filepicker.NewModel(config.InputPath, config.Recursive)
	program := tea.NewProgram(model)

	finalModel, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("TUI error: %w", err)
	}

	// Get the selected file
	if m, ok := finalModel.(filepicker.Model); ok {
		selectedFile := m.GetSelectedFile()
		if selectedFile == "" {
			return "", nil // User cancelled, not an error
		}
		return selectedFile, nil
	}

	return "", fmt.Errorf("unexpected model type")
}
