package main

import (
	"context"
	"testing"

	"github.com/Logta/SurveyBot/types"
)

func TestSurveyWorkflow_Integration(t *testing.T) {
	t.Run("正常系: 完全なアンケートワークフロー", func(t *testing.T) {
		// Arrange
		helper := NewTestHelper(t)
		ctx := context.Background()
		guildID := "integration-test-guild"
		
		surveyHandler := helper.CreateSurveyHandler()

		// Act & Assert: アンケート開始
		if !surveyHandler.CanHandle("!survey") {
			t.Error("Survey handler should handle !survey command")
		}

		// 初期状態の確認
		state, err := helper.StateManager.GetState(ctx, guildID)
		if err != nil {
			t.Fatalf("Failed to get initial state: %v", err)
		}
		if state.Active {
			t.Error("Initial state should not be active")
		}

		// アンケート開始をシミュレート
		helper.SetupSurveyState(t, guildID, true, "")
		helper.VerifySurveyState(t, guildID, true, "")

		// タイトル設定をシミュレート
		helper.SetupSurveyState(t, guildID, true, "統合テストアンケート")
		helper.VerifySurveyState(t, guildID, true, "統合テストアンケート")

		// アンケートキャンセルをシミュレート
		err = helper.StateManager.ClearState(ctx, guildID)
		if err != nil {
			t.Fatalf("Failed to clear state: %v", err)
		}
		helper.VerifySurveyState(t, guildID, false, "")
	})
}

func TestHandlerChain_Integration(t *testing.T) {
	t.Run("正常系: 複数ハンドラーの協調動作", func(t *testing.T) {
		// Arrange
		helper := NewTestHelper(t)
		
		handlers := []types.Handler{
			helper.CreateSurveyHandler(),
			helper.CreateShuffleHandler(),
			helper.CreateCouplingHandler(),
			helper.CreateHelpHandler(),
		}

		testCases := []struct {
			command        string
			expectedHandler string
		}{
			{"!survey", "SurveyHandler"},
			{"!title テスト", "SurveyHandler"},
			{"!shuffle", "ShuffleHandler"},
			{"!coupling", "CouplingHandler"},
			{"!help", "HelpHandler"},
		}

		for _, tc := range testCases {
			t.Run(tc.command, func(t *testing.T) {
				// Act
				var handlingHandler types.Handler
				for _, handler := range handlers {
					if handler.CanHandle(tc.command) {
						handlingHandler = handler
						break
					}
				}

				// Assert
				if handlingHandler == nil {
					t.Errorf("No handler found for command: %s", tc.command)
				} else if handlingHandler.Name() != tc.expectedHandler {
					t.Errorf("Wrong handler for command %s: got %s, want %s", 
						tc.command, handlingHandler.Name(), tc.expectedHandler)
				}
			})
		}
	})
}

func TestUtilsIntegration(t *testing.T) {
	t.Run("正常系: EmojiProvider, Shuffler, Coupler の統合", func(t *testing.T) {
		// Arrange
		helper := NewTestHelper(t)
		ctx := context.Background()

		// EmojiProvider のテスト
		emojis := make([]string, helper.EmojiProvider.GetMaxEmojis())
		for i := 0; i < helper.EmojiProvider.GetMaxEmojis(); i++ {
			emoji, err := helper.EmojiProvider.GetEmoji(ctx, i)
			if err != nil {
				t.Errorf("Failed to get emoji at index %d: %v", i, err)
			}
			emojis[i] = emoji
		}

		// 重複がないことを確認
		emojiSet := make(map[string]bool)
		for _, emoji := range emojis {
			if emojiSet[emoji] {
				t.Errorf("Duplicate emoji found: %s", emoji)
			}
			emojiSet[emoji] = true
		}

		// Shuffler のテスト
		originalItems := []string{"A", "B", "C", "D", "E"}
		shuffledItems := helper.Shuffler.Shuffle(ctx, originalItems)

		if len(shuffledItems) != len(originalItems) {
			t.Errorf("Shuffled items length mismatch: got %d, want %d", 
				len(shuffledItems), len(originalItems))
		}

		// すべてのアイテムが保持されているかを確認
		itemCount := make(map[string]int)
		for _, item := range originalItems {
			itemCount[item]++
		}
		for _, item := range shuffledItems {
			itemCount[item]--
		}
		for item, count := range itemCount {
			if count != 0 {
				t.Errorf("Item count mismatch for %s: %d", item, count)
			}
		}

		// Coupler のテスト
		itemSets := [][]string{
			{"Red", "Blue"},
			{"Circle", "Square"},
		}
		
		couples, err := helper.Coupler.Couple(ctx, itemSets)
		if err != nil {
			t.Errorf("Failed to couple items: %v", err)
		}

		if len(couples) != 2 { // Red,Blue それぞれに対して
			t.Errorf("Unexpected number of couples: got %d, want 2", len(couples))
		}

		for i, couple := range couples {
			if len(couple) != 2 {
				t.Errorf("Couple %d has wrong length: got %d, want 2", i, len(couple))
			}
		}
	})
}

func TestStateManagement_Integration(t *testing.T) {
	t.Run("正常系: 複数ギルドの状態管理", func(t *testing.T) {
		// Arrange
		helper := NewTestHelper(t)
		ctx := context.Background()
		
		guilds := []string{"guild1", "guild2", "guild3"}
		titles := []string{"アンケート1", "アンケート2", "アンケート3"}

		// Act: 各ギルドに異なる状態を設定
		for i, guildID := range guilds {
			helper.SetupSurveyState(t, guildID, true, titles[i])
		}

		// Assert: 各ギルドの状態が独立していることを確認
		for i, guildID := range guilds {
			helper.VerifySurveyState(t, guildID, true, titles[i])
		}

		// Act: 一つのギルドの状態をクリア
		err := helper.StateManager.ClearState(ctx, guilds[1])
		if err != nil {
			t.Fatalf("Failed to clear state for guild2: %v", err)
		}

		// Assert: 他のギルドの状態が影響されていないことを確認
		helper.VerifySurveyState(t, guilds[0], true, titles[0])
		helper.VerifySurveyState(t, guilds[1], false, "") // クリアされた
		helper.VerifySurveyState(t, guilds[2], true, titles[2])
	})
}

func TestErrorHandling_Integration(t *testing.T) {
	t.Run("正常系: エラー耐性の確認", func(t *testing.T) {
		// Arrange
		helper := NewTestHelper(t)
		ctx := context.Background()

		// Act & Assert: 存在しないギルドの状態取得
		state, err := helper.StateManager.GetState(ctx, "non-existent-guild")
		if err != nil {
			t.Errorf("Getting non-existent guild state should not error: %v", err)
		}
		if state.Active {
			t.Error("Non-existent guild should have inactive state")
		}

		// Act & Assert: 範囲外の絵文字取得
		_, err = helper.EmojiProvider.GetEmoji(ctx, 999)
		if err == nil {
			t.Error("Getting out-of-range emoji should return error")
		}

		// Act & Assert: 空の配列のシャッフル
		result := helper.Shuffler.Shuffle(ctx, []string{})
		if len(result) != 0 {
			t.Errorf("Shuffling empty array should return empty array, got length %d", len(result))
		}

		// Act & Assert: 空の配列のカップリング
		_, err = helper.Coupler.Couple(ctx, [][]string{})
		if err == nil {
			t.Error("Coupling empty array should return error")
		}
	})
}