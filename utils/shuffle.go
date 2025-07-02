package utils

import (
	"context"
	"math/rand/v2"

	"github.com/Logta/SurveyBot/types"
)

type shuffler struct{}

// NewShuffler creates a new shuffler instance
func NewShuffler() types.Shuffler {
	return &shuffler{}
}

func (s *shuffler) Shuffle(ctx context.Context, items []string) []string {
	if len(items) <= 1 {
		return items
	}

	// Create a copy to avoid modifying the original slice
	result := make([]string, len(items))
	copy(result, items)

	// Fisher-Yates shuffle
	n := len(result)
	for i := n - 1; i >= 0; i-- {
		j := rand.IntN(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// Legacy function for backward compatibility
func FisherYatesShuffle(data []string) []string {
	shuffler := NewShuffler()
	return shuffler.Shuffle(context.Background(), data)
}
