package utils

import (
	"context"
	"reflect"
	"testing"
)

func TestShuffler_Shuffle(t *testing.T) {
	t.Run("正常系: 複数要素のシャッフル", func(t *testing.T) {
		// Arrange
		shuffler := NewShuffler()
		ctx := context.Background()
		input := []string{"apple", "banana", "cherry", "date", "elderberry"}
		originalInput := make([]string, len(input))
		copy(originalInput, input)

		// Act
		result := shuffler.Shuffle(ctx, input)

		// Assert
		if len(result) != len(originalInput) {
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), len(originalInput))
		}

		// 元の配列が変更されていないことを確認
		if !reflect.DeepEqual(input, originalInput) {
			t.Error("元の配列が変更されてしまいました")
		}

		// すべての要素が含まれていることを確認
		for _, original := range originalInput {
			found := false
			for _, shuffled := range result {
				if original == shuffled {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("要素 %v がシャッフル結果に含まれていません", original)
			}
		}

		// 重複がないことを確認
		seen := make(map[string]bool)
		for _, item := range result {
			if seen[item] {
				t.Errorf("重複した要素が見つかりました: %v", item)
			}
			seen[item] = true
		}
	})

	t.Run("正常系: 空配列", func(t *testing.T) {
		// Arrange
		shuffler := NewShuffler()
		ctx := context.Background()
		input := []string{}

		// Act
		result := shuffler.Shuffle(ctx, input)

		// Assert
		if len(result) != 0 {
			t.Errorf("空配列の結果が期待値と異なります: got %v, want %v", len(result), 0)
		}
	})

	t.Run("正常系: 単一要素", func(t *testing.T) {
		// Arrange
		shuffler := NewShuffler()
		ctx := context.Background()
		input := []string{"single"}

		// Act
		result := shuffler.Shuffle(ctx, input)

		// Assert
		if len(result) != 1 {
			t.Errorf("単一要素の結果長が期待値と異なります: got %v, want %v", len(result), 1)
		}
		if result[0] != "single" {
			t.Errorf("単一要素の結果が期待値と異なります: got %v, want %v", result[0], "single")
		}
	})

	t.Run("正常系: シャッフルの統計的検証", func(t *testing.T) {
		// Arrange
		shuffler := NewShuffler()
		ctx := context.Background()
		input := []string{"A", "B", "C"}
		iterations := 1000
		positionCounts := make(map[string]map[int]int)

		// 各要素の位置カウントを初期化
		for _, item := range input {
			positionCounts[item] = make(map[int]int)
		}

		// Act
		for i := 0; i < iterations; i++ {
			result := shuffler.Shuffle(ctx, input)
			for pos, item := range result {
				positionCounts[item][pos]++
			}
		}

		// Assert
		// 各要素が各位置にある確率がおおよそ1/3になることを確認
		tolerance := 0.15 // 15%の許容誤差
		expectedFreq := float64(iterations) / float64(len(input))

		for item, positions := range positionCounts {
			for pos, count := range positions {
				freq := float64(count)
				if freq < expectedFreq*(1-tolerance) || freq > expectedFreq*(1+tolerance) {
					t.Errorf("要素 %v の位置 %v の出現頻度が期待範囲外です: got %v, expected around %v",
						item, pos, freq, expectedFreq)
				}
			}
		}
	})
}

func TestFisherYatesShuffle_BackwardCompatibility(t *testing.T) {
	t.Run("正常系: 後方互換性の確認", func(t *testing.T) {
		// Arrange
		input := []string{"test", "shuffle", "function"}
		originalLength := len(input)

		// Act
		result := FisherYatesShuffle(input)

		// Assert
		if len(result) != originalLength {
			t.Errorf("結果の長さが期待値と異なります: got %v, want %v", len(result), originalLength)
		}

		// すべての要素が含まれていることを確認
		originalMap := make(map[string]bool)
		for _, item := range []string{"test", "shuffle", "function"} {
			originalMap[item] = true
		}

		for _, item := range result {
			if !originalMap[item] {
				t.Errorf("期待していない要素が結果に含まれています: %v", item)
			}
		}
	})
}
