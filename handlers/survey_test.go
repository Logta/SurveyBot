package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/Logta/SurveyBot/types"
)

// モックの実装
type mockStateManager struct {
	states map[string]*types.SurveyState
	err    error
}

func (m *mockStateManager) GetState(ctx context.Context, guildID string) (*types.SurveyState, error) {
	if m.err != nil {
		return nil, m.err
	}
	if state, exists := m.states[guildID]; exists {
		return &types.SurveyState{Active: state.Active, Title: state.Title}, nil
	}
	return &types.SurveyState{}, nil
}

func (m *mockStateManager) SetState(ctx context.Context, guildID string, state *types.SurveyState) error {
	if m.err != nil {
		return m.err
	}
	if m.states == nil {
		m.states = make(map[string]*types.SurveyState)
	}
	m.states[guildID] = &types.SurveyState{Active: state.Active, Title: state.Title}
	return nil
}

func (m *mockStateManager) ClearState(ctx context.Context, guildID string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.states, guildID)
	return nil
}

type mockEmojiProvider struct {
	emojis []string
	err    error
}

func (m *mockEmojiProvider) GetEmoji(ctx context.Context, index int) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if index < 0 || index >= len(m.emojis) {
		return "", errors.New("index out of range")
	}
	return m.emojis[index], nil
}

func (m *mockEmojiProvider) GetMaxEmojis() int {
	return len(m.emojis)
}

type mockLogger struct {
	logs []string
}

func (m *mockLogger) Info(ctx context.Context, msg string, fields ...types.Field) {
	m.logs = append(m.logs, "INFO: "+msg)
}

func (m *mockLogger) Error(ctx context.Context, msg string, err error, fields ...types.Field) {
	m.logs = append(m.logs, "ERROR: "+msg+" - "+err.Error())
}

func (m *mockLogger) Debug(ctx context.Context, msg string, fields ...types.Field) {
	m.logs = append(m.logs, "DEBUG: "+msg)
}

func TestSurveyHandler_Name(t *testing.T) {
	t.Run("正常系: ハンドラー名を取得", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		// Act
		name := handler.Name()

		// Assert
		expected := "SurveyHandler"
		if name != expected {
			t.Errorf("ハンドラー名が期待値と異なります: got %v, want %v", name, expected)
		}
	})
}

func TestSurveyHandler_CanHandle(t *testing.T) {
	t.Run("正常系: 対応可能なコマンドの判定", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		testCases := []struct {
			command  string
			expected bool
		}{
			{"!survey", true},
			{"!title", true},
			{"!title テストタイトル", true},
			{"!content", true},
			{"!content 選択肢1", true},
			{"!cancel", true},
			{"!check state", true},
			{"!check title", true},
			{"!help", false},
			{"!shuffle", false},
			{"hello", false},
		}

		for _, tc := range testCases {
			t.Run(tc.command, func(t *testing.T) {
				// Act
				result := handler.CanHandle(tc.command)

				// Assert
				if result != tc.expected {
					t.Errorf("コマンド判定が期待値と異なります: command=%v, got=%v, want=%v",
						tc.command, result, tc.expected)
				}
			})
		}
	})
}

func TestSurveyHandler_StateManagement(t *testing.T) {
	t.Run("正常系: 状態管理の基本動作", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{states: make(map[string]*types.SurveyState)}
		ctx := context.Background()
		guildID := "test-guild"

		// Act: 状態を設定
		state := &types.SurveyState{Active: true, Title: "テストタイトル"}
		err := stateManager.SetState(ctx, guildID, state)

		// Assert: 設定が成功
		if err != nil {
			t.Errorf("状態設定でエラーが発生: %v", err)
		}

		// Act: 状態を取得
		retrievedState, err := stateManager.GetState(ctx, guildID)

		// Assert: 正しい状態が取得できる
		if err != nil {
			t.Errorf("状態取得でエラーが発生: %v", err)
		}
		if !retrievedState.Active {
			t.Error("状態がアクティブになっていません")
		}
		if retrievedState.Title != "テストタイトル" {
			t.Errorf("タイトルが期待値と異なります: got %v, want %v", retrievedState.Title, "テストタイトル")
		}

		// Act: 状態をクリア
		err = stateManager.ClearState(ctx, guildID)

		// Assert: クリアが成功
		if err != nil {
			t.Errorf("状態クリアでエラーが発生: %v", err)
		}

		// Act: クリア後の状態を確認
		clearedState, err := stateManager.GetState(ctx, guildID)

		// Assert: 状態がクリアされている
		if err != nil {
			t.Errorf("状態取得でエラーが発生: %v", err)
		}
		if clearedState.Active {
			t.Error("状態がクリアされていません")
		}
		if clearedState.Title != "" {
			t.Error("タイトルがクリアされていません")
		}
	})

	t.Run("異常系: エラーハンドリング", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{err: errors.New("state error")}
		ctx := context.Background()
		guildID := "test-guild"

		// Act & Assert: 設定エラー
		state := &types.SurveyState{Active: true, Title: "テスト"}
		err := stateManager.SetState(ctx, guildID, state)
		if err == nil {
			t.Error("設定エラーが期待されていましたが、nilが返されました")
		}

		// Act & Assert: 取得エラー
		_, err = stateManager.GetState(ctx, guildID)
		if err == nil {
			t.Error("取得エラーが期待されていましたが、nilが返されました")
		}

		// Act & Assert: クリアエラー
		err = stateManager.ClearState(ctx, guildID)
		if err == nil {
			t.Error("クリアエラーが期待されていましたが、nilが返されました")
		}
	})
}
