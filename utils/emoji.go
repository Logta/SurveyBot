package utils

import (
	"context"
	"fmt"

	"github.com/Logta/SurveyBot/types"
)

type emojiProvider struct {
	emojis []string
}

// NewEmojiProvider creates a new emoji provider
func NewEmojiProvider() types.EmojiProvider {
	return &emojiProvider{
		emojis: []string{
			"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣",
			"5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟",
		},
	}
}

func (e *emojiProvider) GetEmoji(ctx context.Context, index int) (string, error) {
	if index < 0 || index >= len(e.emojis) {
		return "", fmt.Errorf("emoji index %d out of range [0-%d]", index, len(e.emojis)-1)
	}
	return e.emojis[index], nil
}

func (e *emojiProvider) GetMaxEmojis() int {
	return len(e.emojis)
}

// Legacy function for backward compatibility
func FindEmoji(num int) (string, error) {
	provider := NewEmojiProvider()
	return provider.GetEmoji(context.Background(), num)
}
