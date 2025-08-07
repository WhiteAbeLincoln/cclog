package formatter

import (
	"testing"
	"time"

	"github.com/annenpolka/cclog/internal/domain"
	"github.com/annenpolka/cclog/internal/testutil"
)

func TestIsContentfulMessage(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	tests := []struct {
		name     string
		message  domain.Message
		expected bool
	}{
		{
			name: "normal user message",
			message: domain.Message{
				Type:      "user",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Hello, how are you?",
				},
			},
			expected: true,
		},
		{
			name: "normal assistant message",
			message: domain.Message{
				Type:      "assistant",
				Timestamp: timestamp,
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
			expected: true,
		},
		{
			name: "system message should be filtered",
			message: domain.Message{
				Type:      "system",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "system",
					"content": "System reminder",
				},
			},
			expected: false,
		},
		{
			name: "empty message should be filtered",
			message: domain.Message{
				Type:      "assistant",
				Timestamp: timestamp,
				Message:   map[string]interface{}{},
			},
			expected: false,
		},
		{
			name: "API error message should be filtered",
			message: domain.Message{
				Type:      "assistant",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "assistant",
					"content": "API Error: Request was aborted.",
				},
			},
			expected: false,
		},
		{
			name: "interrupted request should be filtered",
			message: domain.Message{
				Type:      "user",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "[Request interrupted by user]",
				},
			},
			expected: false,
		},
		{
			name: "command message should be filtered",
			message: domain.Message{
				Type:      "user",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "<command-name>/add-dir</command-name>",
				},
			},
			expected: false,
		},
		{
			name: "bash input should be filtered",
			message: domain.Message{
				Type:      "user",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "<bash-input>git status</bash-input>",
				},
			},
			expected: false,
		},
		{
			name: "meta message should be filtered",
			message: domain.Message{
				Type:      "user",
				Timestamp: timestamp,
				IsMeta:    true,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Some meta content",
				},
			},
			expected: false,
		},
		{
			name: "summary message should be filtered",
			message: domain.Message{
				Type:      "summary",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"summary": "Test conversation summary",
				},
			},
			expected: false,
		},
		{
			name: "message with only UUID should be filtered",
			message: domain.Message{
				Type:      "assistant",
				Timestamp: timestamp,
				UUID:      "some-uuid",
				Message:   nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsContentfulMessage(tt.message)
			testutil.Diff(t, tt.expected, result)
		})
	}
}

func TestFilterMessages(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	messages := []domain.Message{
		{
			Type:      "user",
			Timestamp: timestamp,
			Message: map[string]interface{}{
				"role":    "user",
				"content": "Hello",
			},
		},
		{
			Type:      "system",
			Timestamp: timestamp,
			Message: map[string]interface{}{
				"role":    "system",
				"content": "System message",
			},
		},
		{
			Type:      "assistant",
			Timestamp: timestamp,
			Message: map[string]interface{}{
				"role":    "assistant",
				"content": "Response",
			},
		},
		{
			Type:      "assistant",
			Timestamp: timestamp,
			Message:   nil,
		},
	}

	filtered := FilterMessages(messages, true)

	testutil.Diff(t, 2, len(filtered))

	// Test with filtering disabled
	unfiltered := FilterMessages(messages, false)

	testutil.Diff(t, 4, len(unfiltered))
}

func TestFilterConversationLog(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "2025-07-06T05:01:29.618Z")

	log := &domain.ConversationLog{
		FilePath: "/test/path.jsonl",
		Messages: []domain.Message{
			{
				Type:      "user",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "user",
					"content": "Test message",
				},
			},
			{
				Type:      "system",
				Timestamp: timestamp,
				Message: map[string]interface{}{
					"role":    "system",
					"content": "System message",
				},
			},
		},
	}

	filtered := FilterConversationLog(log, true)

	testutil.Diff(t, 1, len(filtered.Messages))

	testutil.Diff(t, log.FilePath, filtered.FilePath)
}
