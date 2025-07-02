package state

import (
	"context"
	"sync"
	"testing"

	"github.com/Logta/SurveyBot/types"
)

func TestMemoryStateManager_GetState(t *testing.T) {
	t.Run("正常系: 存在するギルドの状態を取得", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-123"
		expectedState := &types.SurveyState{
			Active: true,
			Title:  "テストアンケート",
		}

		// 事前に状態を設定
		err := manager.SetState(ctx, guildID, expectedState)
		if err != nil {
			t.Fatalf("事前設定でエラーが発生: %v", err)
		}

		// Act
		result, err := manager.GetState(ctx, guildID)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}
		if result.Active != expectedState.Active {
			t.Errorf("Active状態が期待値と異なります: got %v, want %v", result.Active, expectedState.Active)
		}
		if result.Title != expectedState.Title {
			t.Errorf("Titleが期待値と異なります: got %v, want %v", result.Title, expectedState.Title)
		}
	})

	t.Run("正常系: 存在しないギルドの状態を取得", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "non-existent-guild"

		// Act
		result, err := manager.GetState(ctx, guildID)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}
		if result.Active != false {
			t.Errorf("初期状態のActiveが期待値と異なります: got %v, want %v", result.Active, false)
		}
		if result.Title != "" {
			t.Errorf("初期状態のTitleが期待値と異なります: got %v, want %v", result.Title, "")
		}
	})

	t.Run("正常系: 状態のコピーが返されることを確認", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-copy"
		originalState := &types.SurveyState{
			Active: true,
			Title:  "オリジナル",
		}

		err := manager.SetState(ctx, guildID, originalState)
		if err != nil {
			t.Fatalf("事前設定でエラーが発生: %v", err)
		}

		// Act
		result1, _ := manager.GetState(ctx, guildID)
		result2, _ := manager.GetState(ctx, guildID)

		// 返された状態を変更
		result1.Title = "変更後"

		// Assert
		if result2.Title != "オリジナル" {
			t.Error("状態のコピーが正しく動作していません。元の状態が変更されています")
		}
	})
}

func TestMemoryStateManager_SetState(t *testing.T) {
	t.Run("正常系: 新しい状態を設定", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-set"
		newState := &types.SurveyState{
			Active: true,
			Title:  "新しいアンケート",
		}

		// Act
		err := manager.SetState(ctx, guildID, newState)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// 設定された状態を確認
		result, _ := manager.GetState(ctx, guildID)
		if result.Active != newState.Active {
			t.Errorf("設定されたActive状態が期待値と異なります: got %v, want %v", result.Active, newState.Active)
		}
		if result.Title != newState.Title {
			t.Errorf("設定されたTitleが期待値と異なります: got %v, want %v", result.Title, newState.Title)
		}
	})

	t.Run("正常系: 既存の状態を更新", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-update"

		initialState := &types.SurveyState{
			Active: false,
			Title:  "初期タイトル",
		}
		updatedState := &types.SurveyState{
			Active: true,
			Title:  "更新後タイトル",
		}

		manager.SetState(ctx, guildID, initialState)

		// Act
		err := manager.SetState(ctx, guildID, updatedState)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		result, _ := manager.GetState(ctx, guildID)
		if result.Active != updatedState.Active {
			t.Errorf("更新されたActive状態が期待値と異なります: got %v, want %v", result.Active, updatedState.Active)
		}
		if result.Title != updatedState.Title {
			t.Errorf("更新されたTitleが期待値と異なります: got %v, want %v", result.Title, updatedState.Title)
		}
	})

	t.Run("異常系: nilの状態を設定", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-nil"

		// Act
		err := manager.SetState(ctx, guildID, nil)

		// Assert
		if err == nil {
			t.Error("エラーが期待されていましたが、nilが返されました")
		}
	})

	t.Run("正常系: 状態のコピーが保存されることを確認", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-copy-set"
		originalState := &types.SurveyState{
			Active: true,
			Title:  "元のタイトル",
		}

		// Act
		err := manager.SetState(ctx, guildID, originalState)
		if err != nil {
			t.Fatalf("SetStateでエラーが発生: %v", err)
		}

		// 元の状態を変更
		originalState.Title = "変更後のタイトル"

		// Assert
		result, _ := manager.GetState(ctx, guildID)
		if result.Title != "元のタイトル" {
			t.Error("状態のコピーが正しく保存されていません。元のオブジェクトの変更が影響しています")
		}
	})
}

func TestMemoryStateManager_ClearState(t *testing.T) {
	t.Run("正常系: 存在する状態をクリア", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-clear"
		state := &types.SurveyState{
			Active: true,
			Title:  "削除予定",
		}

		manager.SetState(ctx, guildID, state)

		// Act
		err := manager.ClearState(ctx, guildID)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// 状態がクリアされていることを確認
		result, _ := manager.GetState(ctx, guildID)
		if result.Active != false {
			t.Errorf("クリア後のActive状態が期待値と異なります: got %v, want %v", result.Active, false)
		}
		if result.Title != "" {
			t.Errorf("クリア後のTitleが期待値と異なります: got %v, want %v", result.Title, "")
		}
	})

	t.Run("正常系: 存在しない状態をクリア", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "non-existent-guild-clear"

		// Act
		err := manager.ClearState(ctx, guildID)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}
	})
}

func TestMemoryStateManager_Concurrency(t *testing.T) {
	t.Run("正常系: 並行アクセスの安全性", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()
		guildID := "test-guild-concurrent"

		var wg sync.WaitGroup
		numGoroutines := 100
		numOperations := 10

		// Act
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < numOperations; j++ {
					// 並行して設定・取得・クリアを実行
					state := &types.SurveyState{
						Active: true,
						Title:  "concurrent test",
					}

					manager.SetState(ctx, guildID, state)
					manager.GetState(ctx, guildID)
					manager.ClearState(ctx, guildID)
				}
			}(i)
		}

		// Assert
		// デッドロックやpanicが発生しないことを確認
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// 正常終了
		case <-ctx.Done():
			t.Error("並行アクセステストがタイムアウトしました")
		}
	})

	t.Run("正常系: 複数ギルドの並行操作", func(t *testing.T) {
		// Arrange
		manager := NewMemoryStateManager()
		ctx := context.Background()

		var wg sync.WaitGroup
		numGuilds := 50

		// Act
		for i := 0; i < numGuilds; i++ {
			wg.Add(1)
			go func(guildNum int) {
				defer wg.Done()

				guildID := string(rune('A' + guildNum))
				state := &types.SurveyState{
					Active: true,
					Title:  guildID + "のアンケート",
				}

				manager.SetState(ctx, guildID, state)
				result, _ := manager.GetState(ctx, guildID)

				if result.Title != state.Title {
					t.Errorf("ギルド %s の状態が期待値と異なります", guildID)
				}
			}(i)
		}

		wg.Wait()

		// Assert
		// 各ギルドの状態が正しく保存されていることを確認
		for i := 0; i < numGuilds; i++ {
			guildID := string(rune('A' + i))
			result, _ := manager.GetState(ctx, guildID)
			expectedTitle := guildID + "のアンケート"

			if result.Title != expectedTitle {
				t.Errorf("ギルド %s の最終状態が期待値と異なります: got %v, want %v",
					guildID, result.Title, expectedTitle)
			}
		}
	})
}
