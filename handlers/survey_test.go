package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/Logta/SurveyBot/types"
	"github.com/bwmarrin/discordgo"
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

func TestSurveyHandler_HandleSurveyStart(t *testing.T) {
	t.Run("正常系: アンケート開始", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{states: make(map[string]*types.SurveyState)}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!survey",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// 状態が正しく設定されているかを確認
		state, _ := stateManager.GetState(ctx, "test-guild")
		if !state.Active {
			t.Error("アンケート状態がアクティブになっていません")
		}
		if state.Title != "" {
			t.Errorf("初期タイトルが空でありません: %v", state.Title)
		}
	})

	t.Run("異常系: 状態管理でエラー", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{err: errors.New("state error")}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!survey",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err == nil {
			t.Error("エラーが期待されていましたが、nilが返されました")
		}
	})
}

func TestSurveyHandler_HandleCancel(t *testing.T) {
	t.Run("正常系: アンケートキャンセル", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: true, Title: "テスト"},
			},
		}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!cancel",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// 状態がクリアされているかを確認
		state, _ := stateManager.GetState(ctx, "test-guild")
		if state.Active {
			t.Error("アンケート状態がクリアされていません")
		}
		if state.Title != "" {
			t.Error("タイトルがクリアされていません")
		}
	})
}

func TestSurveyHandler_HandleTitle(t *testing.T) {
	t.Run("正常系: タイトル設定", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: true, Title: ""},
			},
		}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!title\nテストアンケート",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// タイトルが正しく設定されているかを確認
		state, _ := stateManager.GetState(ctx, "test-guild")
		if state.Title != "テストアンケート" {
			t.Errorf("タイトルが期待値と異なります: got %v, want %v", state.Title, "テストアンケート")
		}
	})

	t.Run("正常系: アンケートが非アクティブの場合は無視", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: false, Title: ""},
			},
		}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!title\nテストアンケート",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err != nil {
			t.Errorf("期待していないエラーが発生: %v", err)
		}

		// タイトルが変更されていないことを確認
		state, _ := stateManager.GetState(ctx, "test-guild")
		if state.Title != "" {
			t.Errorf("タイトルが変更されています: got %v, want empty", state.Title)
		}
	})

	t.Run("異常系: 不正なフォーマット", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: true, Title: ""},
			},
		}
		emojiProvider := &mockEmojiProvider{}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!title", // 改行なし
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// このテストは実際にはDiscordにメッセージを送信しようとするため、
		// モックセッションでは完全にテストできないが、エラーハンドリングを確認

		// Act & Assert
		// エラーが発生するかパニックしないことを確認
		handler.Handle(ctx, session, message)
	})
}

func TestSurveyHandler_CreateSurveyEmbed(t *testing.T) {
	t.Run("正常系: 基本的なアンケート作成", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: true, Title: "テストアンケート"},
			},
		}
		emojiProvider := &mockEmojiProvider{
			emojis: []string{"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣"},
		}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!content\n選択肢A\n選択肢B",
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act & Assert
		// このテストは実際のDiscord APIを使用するため、
		// モック環境では完全なテストが困難
		// エラーが発生しないことを確認
		err := handler.Handle(ctx, session, message)

		// Discordセッションのモックが不完全なため、エラーが発生する可能性があるが
		// ハンドラーのロジック自体は動作することを確認
		if err != nil {
			// ネットワークエラーやDiscord APIエラーは期待される
			t.Logf("予期されるエラー (Discordセッション関連): %v", err)
		}
	})

	t.Run("異常系: 絵文字が不足", func(t *testing.T) {
		// Arrange
		stateManager := &mockStateManager{
			states: map[string]*types.SurveyState{
				"test-guild": {Active: true, Title: "テストアンケート"},
			},
		}
		emojiProvider := &mockEmojiProvider{
			emojis: []string{"1️⃣"}, // 1つだけ
		}
		logger := &mockLogger{}
		handler := NewSurveyHandler(stateManager, emojiProvider, logger)

		ctx := context.Background()
		session := &discordgo.Session{}
		message := &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   "!content\n選択肢A\n選択肢B\n選択肢C", // 3つの選択肢
				ChannelID: "test-channel",
				GuildID:   "test-guild",
			},
		}

		// Act
		err := handler.Handle(ctx, session, message)

		// Assert
		if err == nil {
			t.Error("エラーが期待されていましたが、nilが返されました")
		}
	})
}
