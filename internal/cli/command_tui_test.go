package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/annenpolka/cclog/internal/testutil"
)

func TestParseArgs_TUIMode(t *testing.T) {
	args := []string{"cclog", "--tui"}
	config, err := ParseArgs(args)

	testutil.Diff(t, false, err != nil)
	testutil.True(t, config.TUIMode)
}

func TestParseArgs_TUIModeWithDirectory(t *testing.T) {
	args := []string{"cclog", "--tui", "/path/to/logs"}
	config, err := ParseArgs(args)

	testutil.Diff(t, false, err != nil)
	testutil.True(t, config.TUIMode)
	testutil.Diff(t, "/path/to/logs", config.InputPath)
}

func TestParseArgs_TUIModeDefaultDirectory(t *testing.T) {
	args := []string{"cclog", "--tui"}
	config, err := ParseArgs(args)

	testutil.Diff(t, false, err != nil)
	// Check that config.InputPath is set to either the default directory or fallback
	expectedDir := getDefaultTUIDirectory()
	testutil.True(t, config.InputPath == expectedDir || config.InputPath == ".")
}

func TestRunCommand_TUIMode(t *testing.T) {
	config := Config{
		TUIMode:   true,
		InputPath: ".",
	}

	// This should return empty string and no error for TUI mode
	// since TUI mode is handled differently
	output, err := RunCommand(config)

	testutil.Diff(t, false, err != nil)
	testutil.Diff(t, "", output)
}

func TestEnsureDefaultDirectoryExists(t *testing.T) {
	// Create a temporary directory to simulate user home
	tempHome := t.TempDir()
	testDir := filepath.Join(tempHome, ".claude", "projects")

	// Directory should not exist initially
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Errorf("Test directory should not exist initially")
	}

	// Call the function to check if directory exists
	err := ensureDefaultDirectoryExists(testDir)
	testutil.True(t, err != nil)

	// Create directory first
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Now the function should return no error
	err = ensureDefaultDirectoryExists(testDir)
	testutil.Diff(t, false, err != nil)
}

func TestEnsureDefaultDirectoryExists_AlreadyExists(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "existing")

	// Create the directory first
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Call the function - should not error on existing directory
	err = ensureDefaultDirectoryExists(testDir)
	testutil.Diff(t, false, err != nil)

	// Directory should still exist
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Directory should still exist")
	}
}
