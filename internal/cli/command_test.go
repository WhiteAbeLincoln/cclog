package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/annenpolka/cclog/internal/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
		wantErr  bool
	}{
		{
			name: "single file input",
			args: []string{"cclog", "/path/to/file.jsonl"},
			expected: Config{
				InputPath:   "/path/to/file.jsonl",
				OutputPath:  "",
				IsDirectory: false,
			},
			wantErr: false,
		},
		{
			name: "directory input",
			args: []string{"cclog", "-d", "/path/to/dir"},
			expected: Config{
				InputPath:   "/path/to/dir",
				OutputPath:  "",
				IsDirectory: true,
			},
			wantErr: false,
		},
		{
			name: "file with output",
			args: []string{"cclog", "/path/to/file.jsonl", "-o", "output.md"},
			expected: Config{
				InputPath:   "/path/to/file.jsonl",
				OutputPath:  "output.md",
				IsDirectory: false,
			},
			wantErr: false,
		},
		{
			name: "no arguments - should enable TUI mode and recursive mode by default",
			args: []string{"cclog"},
			expected: Config{
				TUIMode:   true,
				Recursive: true,
			},
			wantErr: false,
		},
		{
			name: "help flag",
			args: []string{"cclog", "-h"},
			expected: Config{
				ShowHelp: true,
			},
			wantErr: false,
		},
		{
			name: "include all flag",
			args: []string{"cclog", "/path/to/file.jsonl", "--include-all"},
			expected: Config{
				InputPath:   "/path/to/file.jsonl",
				OutputPath:  "",
				IsDirectory: false,
				IncludeAll:  true,
			},
			wantErr: false,
		},
		{
			name: "explicit TUI mode",
			args: []string{"cclog", "--tui"},
			expected: Config{
				TUIMode: true,
			},
			wantErr: false,
		},
		{
			name: "TUI mode with path",
			args: []string{"cclog", "--tui", "/path/to/logs"},
			expected: Config{
				InputPath: "/path/to/logs",
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "recursive flag",
			args: []string{"cclog", "--recursive", "/path/to/logs"},
			expected: Config{
				InputPath: "/path/to/logs",
				Recursive: true,
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "recursive and TUI mode combined",
			args: []string{"cclog", "--recursive", "--tui"},
			expected: Config{
				Recursive: true,
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "recursive flag alone should enable TUI mode",
			args: []string{"cclog", "--recursive"},
			expected: Config{
				Recursive: true,
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "short recursive flag alone should enable TUI mode",
			args: []string{"cclog", "-r"},
			expected: Config{
				Recursive: true,
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "recursive with path should enable TUI mode",
			args: []string{"cclog", "--recursive", "/path/to/logs"},
			expected: Config{
				InputPath: "/path/to/logs",
				Recursive: true,
				TUIMode:   true,
			},
			wantErr: false,
		},
		{
			name: "path option should set input path",
			args: []string{"cclog", "--path", "/custom/path"},
			expected: Config{
				InputPath: "/custom/path",
				TUIMode:   true,
				Recursive: true,
			},
			wantErr: false,
		},
		{
			name:     "path option without value should return error",
			args:     []string{"cclog", "--path"},
			expected: Config{},
			wantErr:  true,
		},
		{
			name: "proj flag without value should enable TUI and recursive",
			args: []string{"cclog", "--proj"},
			expected: Config{
				TUIMode:   true,
				Recursive: true,
			},
			wantErr: false,
		},
		{
			name:     "proj flag with non-existent project path should return error",
			args:     []string{"cclog", "--proj", "/nonexistent/path"},
			expected: Config{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseArgs(tt.args)

			gotErr := err != nil
			testutil.Diff(t, tt.wantErr, gotErr)

			if !tt.wantErr {
				opts := []cmp.Option{}
				// For TUI mode tests, we don't check InputPath if expected is empty
				if tt.expected.InputPath == "" {
					opts = append(opts, cmpopts.IgnoreFields(Config{}, "InputPath"))
				}
				testutil.Diff(t, tt.expected, config, opts...)
			}
		})
	}
}

func TestRunCommand(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.jsonl")

	testContent := `{"type":"user","message":{"role":"user","content":"test"},"timestamp":"2025-07-06T05:01:29.618Z","uuid":"test-uuid"}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"response"}]},"timestamp":"2025-07-06T05:01:30.618Z","uuid":"test-uuid-2"}`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		InputPath:   testFile,
		OutputPath:  "",
		IsDirectory: false,
	}

	output, err := RunCommand(config)
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	if !strings.Contains(output, "# Conversation Log") {
		t.Error("Output should contain conversation log header")
	}

	if !strings.Contains(output, "test") {
		t.Error("Output should contain test message content")
	}

	if !strings.Contains(output, "response") {
		t.Error("Output should contain response message content")
	}
}

func TestRunCommandWithDirectory(t *testing.T) {
	// Create a temporary directory with test files
	tempDir := t.TempDir()
	testFile1 := filepath.Join(tempDir, "test1.jsonl")
	testFile2 := filepath.Join(tempDir, "test2.jsonl")

	testContent := `{"type":"user","message":{"role":"user","content":"test1"},"timestamp":"2025-07-06T05:01:29.618Z","uuid":"test-uuid"}`

	err := os.WriteFile(testFile1, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}

	testContent2 := `{"type":"user","message":{"role":"user","content":"test2"},"timestamp":"2025-07-06T05:01:29.618Z","uuid":"test-uuid"}`
	err = os.WriteFile(testFile2, []byte(testContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	config := Config{
		InputPath:   tempDir,
		OutputPath:  "",
		IsDirectory: true,
	}

	output, err := RunCommand(config)
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	if !strings.Contains(output, "# Claude Conversation Logs") {
		t.Error("Output should contain multiple conversations header")
	}

	if !strings.Contains(output, "test1") {
		t.Error("Output should contain content from test1")
	}

	if !strings.Contains(output, "test2") {
		t.Error("Output should contain content from test2")
	}
}

func TestRunCommand_ShowTitle_Single(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "titled.jsonl")

	// First user message content should become title
	testContent := `{"type":"user","message":{"role":"user","content":"hello world"},"timestamp":"2025-07-06T05:01:29.618Z","uuid":"u1"}`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := Config{InputPath: testFile, ShowTitle: true}
	out, err := RunCommand(cfg)
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	if !strings.HasPrefix(out, "# hello world\n\n") {
		t.Fatalf("expected markdown to start with title, got: %q", out[:min(40, len(out))])
	}
}

func TestRunCommand_ShowTitle_Directory(t *testing.T) {
	tempDir := t.TempDir()
	f1 := filepath.Join(tempDir, "a.jsonl")
	f2 := filepath.Join(tempDir, "b.jsonl")

	// Directory: title taken from first log after filtering; use f1
	c1 := `{"type":"user","message":{"role":"user","content":"dir-title"},"timestamp":"2025-07-06T05:01:29.618Z","uuid":"u1"}`
	c2 := `{"type":"user","message":{"role":"user","content":"other"},"timestamp":"2025-07-06T05:01:30.618Z","uuid":"u2"}`

	if err := os.WriteFile(f1, []byte(c1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f2, []byte(c2), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{InputPath: tempDir, IsDirectory: true, ShowTitle: true}
	out, err := RunCommand(cfg)
	if err != nil {
		t.Fatalf("RunCommand failed: %v", err)
	}

	if !strings.HasPrefix(out, "# dir-title\n\n") {
		t.Fatalf("expected directory markdown to start with title, got: %q", out[:min(40, len(out))])
	}
}

// helper: safe substring length for error messages
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestGetDefaultTUIDirectory(t *testing.T) {
	defaultDir := getDefaultTUIDirectory()

	// Should contain either .claude/projects or .config/claude/projects
	hasClaudeProjects := strings.Contains(defaultDir, ".claude/projects")
	hasConfigClaudeProjects := strings.Contains(defaultDir, ".config/claude/projects")

	if !hasClaudeProjects && !hasConfigClaudeProjects {
		t.Errorf("Default directory should contain '.claude/projects' or '.config/claude/projects', got: %s", defaultDir)
	}

	// Should be an absolute path
	testutil.True(t, filepath.IsAbs(defaultDir))
}

func TestGetDefaultTUIDirectory_ValidPath(t *testing.T) {
	defaultDir := getDefaultTUIDirectory()

	// Should be a valid path format
	testutil.Diff(t, false, defaultDir == "")

	// Should end with projects
	testutil.True(t, strings.HasSuffix(defaultDir, "projects"))
}

func TestGetDefaultTUIDirectory_FallbackBehavior(t *testing.T) {
	// Create a temporary directory to simulate user home
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")

	defer func() {
		// Restore original HOME
		os.Setenv("HOME", originalHome)
	}()

	// Test case 1: When .claude directory exists, it should be preferred
	os.Setenv("HOME", tempHome)
	claudeDir := filepath.Join(tempHome, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("Failed to create .claude directory: %v", err)
	}

	result := getDefaultTUIDirectory()
	expected := filepath.Join(tempHome, ".claude", "projects")
	testutil.Diff(t, expected, result)

	// Test case 2: When .claude directory doesn't exist, should fallback to .config/claude
	os.RemoveAll(claudeDir)
	result = getDefaultTUIDirectory()
	expected = filepath.Join(tempHome, ".config", "claude", "projects")
	testutil.Diff(t, expected, result)
}

func TestProjectDirName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/Users/abe/Projects/cclog", "-Users-abe-Projects-cclog"},
		{"/Users/abe/a/b", "-Users-abe-a-b"},
		{"/", "-"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			testutil.Diff(t, tt.expected, projectDirName(tt.input))
		})
	}
}

func TestResolveProjectDirectory(t *testing.T) {
	// Set up a fake HOME with .claude/projects structure
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	projectsDir := filepath.Join(tempHome, ".claude", "projects")

	// Create a project directory that simulates /a/b/c having been a Claude session
	projectPath := filepath.Join(tempHome, "a", "b", "c")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatal(err)
	}
	encodedName := projectDirName(filepath.Join(tempHome, "a", "b", "c"))
	claudeProjDir := filepath.Join(projectsDir, encodedName)
	if err := os.MkdirAll(claudeProjDir, 0755); err != nil {
		t.Fatal(err)
	}

	t.Run("exact match", func(t *testing.T) {
		result, err := resolveProjectDirectory(filepath.Join(tempHome, "a", "b", "c"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.Diff(t, claudeProjDir, result)
	})

	t.Run("subdirectory walks up to find parent", func(t *testing.T) {
		// Create subdirectory d/e under the project
		subDir := filepath.Join(tempHome, "a", "b", "c", "d", "e")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatal(err)
		}

		result, err := resolveProjectDirectory(subDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		testutil.Diff(t, claudeProjDir, result)
	})

	t.Run("no match returns error", func(t *testing.T) {
		noMatchDir := filepath.Join(tempHome, "x", "y")
		if err := os.MkdirAll(noMatchDir, 0755); err != nil {
			t.Fatal(err)
		}

		_, err := resolveProjectDirectory(noMatchDir)
		if err == nil {
			t.Fatal("expected error for unmatched directory, got nil")
		}
		if !strings.Contains(err.Error(), "no Claude project directory found") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func TestParseArgs_ProjResolvesPath(t *testing.T) {
	// Set up a fake HOME with .claude/projects structure
	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	projectsDir := filepath.Join(tempHome, ".claude", "projects")

	// Create a project directory matching CWD
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	encodedName := projectDirName(cwd)
	claudeProjDir := filepath.Join(projectsDir, encodedName)
	if err := os.MkdirAll(claudeProjDir, 0755); err != nil {
		t.Fatal(err)
	}

	// --proj without value should resolve "." to the Claude projects dir for CWD
	config, err := ParseArgs([]string{"cclog", "--proj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	testutil.Diff(t, claudeProjDir, config.InputPath)
	testutil.True(t, config.TUIMode)
	testutil.True(t, config.Recursive)
}
