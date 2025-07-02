package utils

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"

	"github.com/Logta/SurveyBot/types"
)

type coupler struct{}

// NewCoupler creates a new coupler instance
func NewCoupler() types.Coupler {
	return &coupler{}
}

func (c *coupler) Couple(ctx context.Context, itemSets [][]string) ([][]string, error) {
	if len(itemSets) == 0 {
		return nil, fmt.Errorf("no item sets provided")
	}

	// Create deep copy to avoid modifying original data
	sets := make([][]string, len(itemSets))
	for i, set := range itemSets {
		sets[i] = make([]string, len(set))
		copy(sets[i], set)
	}

	var result [][]string
	return c.coupleRecursive(sets, result), nil
}

func (c *coupler) coupleRecursive(sets [][]string, result [][]string) [][]string {
	var couple []string
	allEmpty := true

	for i, set := range sets {
		if len(set) > 0 {
			allEmpty = false
			idx := rand.IntN(len(set))
			couple = append(couple, set[idx])
			// Remove selected item
			sets[i] = append(set[:idx], set[idx+1:]...)
		} else {
			couple = append(couple, "")
		}
	}

	if !allEmpty {
		result = append(result, couple)
		return c.coupleRecursive(sets, result)
	}

	return result
}

// ParseItemSets parses comma-separated strings into 2D slice
func ParseItemSets(data []string, separator string) [][]string {
	var result [][]string
	for _, line := range data {
		if strings.TrimSpace(line) == "" {
			continue
		}
		set := regexp.MustCompile(separator).Split(line, -1)
		// Clean up whitespace
		for i, item := range set {
			set[i] = strings.TrimSpace(item)
		}
		result = append(result, set)
	}
	return result
}

// Legacy functions for backward compatibility
func Coupling(lines [][]string, coupling [][]string) [][]string {
	coupler := NewCoupler()
	result, _ := coupler.Couple(context.Background(), lines)
	return result
}

func GetItemSets(data []string, splitter string) [][]string {
	return ParseItemSets(data, splitter)
}
