package filepicker

import (
	"testing"
	"time"

	"github.com/annenpolka/cclog/internal/domain"
	"github.com/annenpolka/cclog/internal/testutil"
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

	testutil.True(t, SearchInConversation("hello", msgs))
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

	testutil.True(t, SearchInConversation("second", msgs))
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

	testutil.False(t, SearchInConversation("magic-uuid", msgs))
}
