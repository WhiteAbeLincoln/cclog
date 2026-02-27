package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/annenpolka/cclog/internal/domain"
	"github.com/annenpolka/cclog/internal/testutil"
)

func TestFormatConversationToMarkdownWithoutUUID(t *testing.T) {
	// Test default behavior (no UUID)
	timestamp1, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	log := &domain.ConversationLog{
		FilePath: "/test/path/sample.jsonl",
		Messages: []domain.Message{
			{
				Type:      "user",
				UUID:      "user-uuid-1",
				Timestamp: timestamp1,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Hello, how are you?",
				},
			},
		},
	}

	markdown := FormatConversationToMarkdown(log)

	// Check that UUID is NOT included by default
	if strings.Contains(markdown, "UUID:") {
		t.Error("Markdown should not contain UUID by default")
	}

	if !strings.Contains(markdown, "Hello, how are you?") {
		t.Error("Markdown should contain user message content")
	}
}

func TestFormatConversationToMarkdownWithUUID(t *testing.T) {
	// Test with UUID enabled
	timestamp1, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	log := &domain.ConversationLog{
		FilePath: "/test/path/sample.jsonl",
		Messages: []domain.Message{
			{
				Type:      "user",
				UUID:      "user-uuid-1",
				Timestamp: timestamp1,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Hello, how are you?",
				},
			},
		},
	}

	markdown := FormatConversationToMarkdown(log, FormatOptions{ShowUUID: true})

	// Check that UUID IS included when enabled
	if !strings.Contains(markdown, "UUID: user-uuid-1") {
		t.Error("Markdown should contain UUID when enabled")
	}
}

func TestFormatConversationToMarkdown(t *testing.T) {
	// Pin local timezone to UTC so date assertions are stable
	originalLocal := time.Local
	time.Local = time.UTC
	t.Cleanup(func() { time.Local = originalLocal })

	// Create test data
	timestamp1, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")
	timestamp2, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:44.663Z")

	log := &domain.ConversationLog{
		FilePath: "/test/path/sample.jsonl",
		Messages: []domain.Message{
			{
				Type:      "user",
				UUID:      "user-uuid-1",
				Timestamp: timestamp1,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Hello, how are you?",
				},
			},
			{
				Type:      "assistant",
				UUID:      "assistant-uuid-1",
				Timestamp: timestamp2,
				Message: map[string]interface{}{
					"role": "assistant",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "I'm doing well, thank you!",
						},
					},
				},
			},
		},
	}

	markdown := FormatConversationToMarkdown(log)

	// Check if markdown contains expected elements
	if !strings.Contains(markdown, "# Conversation Log") {
		t.Error("Markdown should contain main title")
	}

	if !strings.Contains(markdown, "**File:** `/test/path/sample.jsonl`") {
		t.Error("Markdown should contain file path")
	}

	if !strings.Contains(markdown, "## User") {
		t.Error("Markdown should contain user section")
	}

	if !strings.Contains(markdown, "## Assistant") {
		t.Error("Markdown should contain assistant section")
	}

	if !strings.Contains(markdown, "Hello, how are you?") {
		t.Error("Markdown should contain user message content")
	}

	if !strings.Contains(markdown, "I'm doing well, thank you!") {
		t.Error("Markdown should contain assistant message content")
	}

	// Check that timestamp is formatted correctly (depends on system timezone)
	if !strings.Contains(markdown, "2025-07-06") {
		t.Error("Markdown should contain formatted date")
	}

	// Check that timestamp format is correct (HH:MM:SS format)
	if !strings.Contains(markdown, "**Time:**") {
		t.Error("Markdown should contain timestamp label")
	}
}

func TestFormatMultipleConversationsToMarkdown(t *testing.T) {
	timestamp1, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	logs := []*domain.ConversationLog{
		{
			FilePath: "/test/log1.jsonl",
			Messages: []domain.Message{
				{
					Type:      "user",
					UUID:      "user-uuid-1",
					Timestamp: timestamp1,
					Message: map[string]interface{}{
						"role":    "user",
						"content": "First conversation",
					},
				},
			},
		},
		{
			FilePath: "/test/log2.jsonl",
			Messages: []domain.Message{
				{
					Type:      "user",
					UUID:      "user-uuid-2",
					Timestamp: timestamp1,
					Message: map[string]interface{}{
						"role":    "user",
						"content": "Second conversation",
					},
				},
			},
		},
	}

	markdown := FormatMultipleConversationsToMarkdown(logs)

	if !strings.Contains(markdown, "# Claude Conversation Logs") {
		t.Error("Markdown should contain main title for multiple conversations")
	}

	if !strings.Contains(markdown, "First conversation") {
		t.Error("Markdown should contain first conversation content")
	}

	if !strings.Contains(markdown, "Second conversation") {
		t.Error("Markdown should contain second conversation content")
	}

	if !strings.Contains(markdown, "log1.jsonl") {
		t.Error("Markdown should contain first log filename")
	}

	if !strings.Contains(markdown, "log2.jsonl") {
		t.Error("Markdown should contain second log filename")
	}
}

func TestExtractMessageContent(t *testing.T) {
	tests := []struct {
		name     string
		message  interface{}
		expected string
	}{
		{
			name: "simple string content",
			message: map[string]interface{}{
				"role":    "user",
				"content": "Hello world",
			},
			expected: "Hello world",
		},
		{
			name: "complex content array",
			message: map[string]interface{}{
				"role": "assistant",
				"content": []interface{}{
					map[string]interface{}{
						"type": "text",
						"text": "Response text",
					},
				},
			},
			expected: "Response text",
		},
		{
			name:     "nil message",
			message:  nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractMessageContent(tt.message)
			testutil.Diff(t, tt.expected, result)
		})
	}
}

func TestExtractMessageContentWithPlaceholders(t *testing.T) {
	tests := []struct {
		name             string
		message          interface{}
		showPlaceholders bool
		expectedWithout  string
		expectedWith     string
	}{
		{
			name: "meta message with isMeta flag",
			message: map[string]interface{}{
				"role":    "user",
				"content": "Caveat: The messages below were generated by the user while running local commands.",
			},
			showPlaceholders: true,
			expectedWithout:  "Caveat: The messages below were generated by the user while running local commands.",
			expectedWith:     "*[System warning message - contains caveats about local commands]*",
		},
		{
			name: "command execution message",
			message: map[string]interface{}{
				"role":    "user",
				"content": "<command-name>/ide</command-name>\n<command-message>ide</command-message>\n<command-args></command-args>",
			},
			showPlaceholders: true,
			expectedWithout:  "<command-name>/ide</command-name>\n<command-message>ide</command-message>\n<command-args></command-args>",
			expectedWith:     "*[Command executed: /ide]*",
		},
		{
			name: "command output message",
			message: map[string]interface{}{
				"role":    "user",
				"content": "<local-command-stdout>Connected to Visual Studio Code.</local-command-stdout>",
			},
			showPlaceholders: true,
			expectedWithout:  "<local-command-stdout>Connected to Visual Studio Code.</local-command-stdout>",
			expectedWith:     "*[Command output: Connected to Visual Studio Code.]*",
		},
		{
			name: "empty content",
			message: map[string]interface{}{
				"role":    "assistant",
				"content": "",
			},
			showPlaceholders: true,
			expectedWithout:  "",
			expectedWith:     "*[Empty message content]*",
		},
		{
			name: "empty content with tool use result",
			message: map[string]interface{}{
				"role":    "user",
				"content": "",
				"toolUseResult": map[string]interface{}{
					"type":     "create",
					"filePath": "/tmp/test.txt",
					"content":  "",
				},
			},
			showPlaceholders: true,
			expectedWithout:  "",
			expectedWith:     "*[File created: /tmp/test.txt (empty)]*",
		},
		{
			name: "empty content with command result",
			message: map[string]interface{}{
				"role":    "user",
				"content": "",
				"toolUseResult": map[string]interface{}{
					"stdout":      "",
					"stderr":      "",
					"interrupted": false,
				},
			},
			showPlaceholders: true,
			expectedWithout:  "",
			expectedWith:     "*[Command executed successfully (no output)]*",
		},
		{
			name: "normal message unchanged",
			message: map[string]interface{}{
				"role":    "user",
				"content": "This is a normal user message",
			},
			showPlaceholders: true,
			expectedWithout:  "This is a normal user message",
			expectedWith:     "This is a normal user message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test without placeholders (current behavior)
			result := ExtractMessageContent(tt.message)
			testutil.Diff(t, tt.expectedWithout, result)

			// Test with placeholders (new behavior)
			result = ExtractMessageContent(tt.message, tt.showPlaceholders)
			testutil.Diff(t, tt.expectedWith, result)
		})
	}
}

func TestFormatConversationWithSubagentLinks(t *testing.T) {
	// Pin timezone for stable output
	originalLocal := time.Local
	time.Local = time.UTC
	t.Cleanup(func() { time.Local = originalLocal })

	ts1, _ := time.Parse(time.RFC3339, "2025-07-06T10:00:00.000Z")
	ts2, _ := time.Parse(time.RFC3339, "2025-07-06T10:00:05.000Z")
	ts3, _ := time.Parse(time.RFC3339, "2025-07-06T10:01:00.000Z")

	log := &domain.ConversationLog{
		FilePath: "/test/conversation.jsonl",
		Messages: []domain.Message{
			{
				Type:      "user",
				UUID:      "msg-1",
				Timestamp: ts1,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Add unit tests",
				},
			},
			{
				Type:      "assistant",
				UUID:      "msg-2",
				Timestamp: ts2,
				Message: map[string]interface{}{
					"role": "assistant",
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": "Let me explore first.",
						},
						map[string]interface{}{
							"type": "tool_use",
							"name": "Task",
							"input": map[string]interface{}{
								"description": "Explore project",
							},
						},
					},
				},
			},
			{
				Type:      "assistant",
				UUID:      "msg-3",
				Timestamp: ts3,
				Message: map[string]interface{}{
					"role":    "assistant",
					"content": "Now writing tests.",
				},
			},
		},
	}

	subagents := []domain.SubagentInfo{
		{
			FilePath:  "/test/session/subagents/agent-abc.jsonl",
			AgentID:   "abc",
			Title:     "Find source files",
			Timestamp: ts2.Add(3 * time.Millisecond), // 3ms after parent message
		},
	}

	markdown := FormatConversationToMarkdown(log, FormatOptions{
		Subagents: subagents,
	})

	// Subagent link should appear after the assistant message at ts2
	testutil.True(t, strings.Contains(markdown, "- Subagent: [Find source files](/test/session/subagents/agent-abc.jsonl)"))

	// The link should appear between the two assistant messages
	exploreIdx := strings.Index(markdown, "Let me explore first.")
	linkIdx := strings.Index(markdown, "Subagent: [Find source files]")
	writingIdx := strings.Index(markdown, "Now writing tests.")

	testutil.True(t, exploreIdx < linkIdx)
	testutil.True(t, linkIdx < writingIdx)
}

func TestAssignSubagentsToMessages(t *testing.T) {
	ts1, _ := time.Parse(time.RFC3339, "2025-07-06T10:00:00.000Z")
	ts2, _ := time.Parse(time.RFC3339, "2025-07-06T10:00:05.000Z")
	ts3, _ := time.Parse(time.RFC3339, "2025-07-06T10:01:00.000Z")

	messages := []domain.Message{
		{Type: "user", Timestamp: ts1},
		{Type: "assistant", Timestamp: ts2},
		{Type: "assistant", Timestamp: ts3},
	}

	subagents := []domain.SubagentInfo{
		{AgentID: "a1", Timestamp: ts2.Add(3 * time.Millisecond)},
		{AgentID: "a2", Timestamp: ts3.Add(2 * time.Millisecond)},
	}

	result := assignSubagentsToMessages(messages, subagents)

	// a1 should match to message index 1 (assistant at ts2)
	testutil.Diff(t, 1, len(result[1]))
	testutil.Diff(t, "a1", result[1][0].AgentID)

	// a2 should match to message index 2 (assistant at ts3)
	testutil.Diff(t, 1, len(result[2]))
	testutil.Diff(t, "a2", result[2][0].AgentID)

	// No subagents on user message
	testutil.Diff(t, 0, len(result[0]))
}
