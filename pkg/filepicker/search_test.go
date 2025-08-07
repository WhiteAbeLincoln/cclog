package filepicker

import (
	"testing"
	"time"

	"github.com/annenpolka/cclog/internal/domain"
)

// Red: SearchInConversation should match text inside Message.content
func TestSearchInConversation_MatchesTextContent(t *testing.T) {
	msgs := []domain.Message{
		{
			Type:      "user",
			Timestamp: time.Now(),
			Message: map[string]any{
				"role":    "user",
				"content": "hello world",
			},
		},
	}

	if ok := SearchInConversation("hello", msgs); !ok {
		t.Fatalf("expected query to match user message content")
	}
}

// Red: SearchInConversation should match text inside array-based content blocks
func TestSearchInConversation_MatchesArrayContent(t *testing.T) {
	msgs := []domain.Message{
		{
			Type:      "assistant",
			Timestamp: time.Now(),
			Message: map[string]any{
				"role": "assistant",
				"content": []any{
					map[string]any{"type": "text", "text": "first"},
					map[string]any{"type": "text", "text": "second block"},
				},
			},
		},
	}

	if ok := SearchInConversation("second", msgs); !ok {
		t.Fatalf("expected query to match text inside array-based content")
	}
}

// Red: Should not match metadata like UUID; only message content
func TestSearchInConversation_DoesNotMatchMetadata(t *testing.T) {
	msgs := []domain.Message{
		{
			Type:      "assistant",
			UUID:      "magic-uuid-12345",
			Timestamp: time.Now(),
			Message: map[string]any{
				"role":    "assistant",
				"content": "just content here",
			},
		},
	}

	if ok := SearchInConversation("magic-uuid", msgs); ok {
		t.Fatalf("did not expect query to match metadata fields like UUID")
	}
}
