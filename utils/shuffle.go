package utils
import (
    "math/rand"
)

func FisherYatesShuffle(data []string) []string {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}

	return data
}
