package utils

import (
	"context"
	"testing"
)

func TestEmojiProvider_GetEmoji(t *testing.T) {
	t.Run("正常系: 有効なインデックスで絵文字を取得", func(t *testing.T) {
		// Arrange
		provider := NewEmojiProvider()
		ctx := context.Background()
		testCases := []struct {
			index    int
			expected string
		}{
			{0, "0️⃣"},
			{1, "1️⃣"},
			{5, "5️⃣"},
			{10, "🔟"},
		}

		for _, tc := range testCases {
			t.Run(tc.expected, func(t *testing.T) {
				// Act
				result, err := provider.GetEmoji(ctx, tc.index)

				// Assert
				if err != nil {
					t.Errorf("期待していないエラーが発生: %v", err)
				}
				if result != tc.expected {
					t.Errorf("絵文字が期待値と異なります: got %v, want %v", result, tc.expected)
				}
			})
		}
	})

	t.Run("異常系: 範囲外のインデックス", func(t *testing.T) {
		// Arrange
		provider := NewEmojiProvider()
		ctx := context.Background()
		invalidIndices := []int{-1, 11, 100}

		for _, index := range invalidIndices {
			t.Run(string(rune(index+48)), func(t *testing.T) {
				// Act
				result, err := provider.GetEmoji(ctx, index)

				// Assert
				if err == nil {
					t.Error("エラーが期待されていましたが、nilが返されました")
				}
				if result != "" {
					t.Errorf("空文字列が期待されていましたが、%vが返されました", result)
				}
			})
		}
	})
}

func TestEmojiProvider_GetMaxEmojis(t *testing.T) {
	t.Run("正常系: 最大絵文字数を取得", func(t *testing.T) {
		// Arrange
		provider := NewEmojiProvider()
		expected := 11 // 0-9 + 🔟

		// Act
		result := provider.GetMaxEmojis()

		// Assert
		if result != expected {
			t.Errorf("最大絵文字数が期待値と異なります: got %v, want %v", result, expected)
		}
	})
}

func TestFindEmoji_BackwardCompatibility(t *testing.T) {
	t.Run("正常系: 後方互換性の確認", func(t *testing.T) {
		// Arrange
		testCases := []struct {
			input    int
			expected string
		}{
			{0, "0️⃣"},
			{1, "1️⃣"},
			{10, "🔟"},
		}

		for _, tc := range testCases {
			t.Run(tc.expected, func(t *testing.T) {
				// Act
				result, err := FindEmoji(tc.input)

				// Assert
				if err != nil {
					t.Errorf("期待していないエラーが発生: %v", err)
				}
				if result != tc.expected {
					t.Errorf("絵文字が期待値と異なります: got %v, want %v", result, tc.expected)
				}
			})
		}
	})

	t.Run("異常系: 範囲外の値", func(t *testing.T) {
		// Arrange
		invalidValue := 15

		// Act
		result, err := FindEmoji(invalidValue)

		// Assert
		if err == nil {
			t.Error("エラーが期待されていましたが、nilが返されました")
		}
		if result != "" {
			t.Errorf("空文字列が期待されていましたが、%vが返されました", result)
		}
	})
}
