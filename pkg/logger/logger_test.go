package logger

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/Logta/SurveyBot/types"
)

func TestLogger_Info(t *testing.T) {
	t.Run("正常系: 基本的な情報ログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "テスト情報メッセージ"

		// Act
		logger.Info(ctx, message)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[INFO]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
	})

	t.Run("正常系: フィールド付き情報ログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "フィールド付きメッセージ"
		fields := []types.Field{
			{Key: "user_id", Value: "12345"},
			{Key: "action", Value: "login"},
		}

		// Act
		logger.Info(ctx, message, fields...)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[INFO]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
		if !strings.Contains(output, "user_id=12345") {
			t.Errorf("フィールドuser_idが含まれていません: %v", output)
		}
		if !strings.Contains(output, "action=login") {
			t.Errorf("フィールドactionが含まれていません: %v", output)
		}
	})
}

func TestLogger_Error(t *testing.T) {
	t.Run("正常系: 基本的なエラーログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "エラーが発生しました"
		err := errors.New("テストエラー")

		// Act
		logger.Error(ctx, message, err)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[ERROR]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
		if !strings.Contains(output, "error=テストエラー") {
			t.Errorf("エラー情報が含まれていません: %v", output)
		}
	})

	t.Run("正常系: フィールド付きエラーログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "データベース接続エラー"
		err := errors.New("connection timeout")
		fields := []types.Field{
			{Key: "database", Value: "users"},
			{Key: "timeout", Value: 30},
		}

		// Act
		logger.Error(ctx, message, err, fields...)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[ERROR]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
		if !strings.Contains(output, "error=connection timeout") {
			t.Errorf("エラー情報が含まれていません: %v", output)
		}
		if !strings.Contains(output, "database=users") {
			t.Errorf("フィールドdatabaseが含まれていません: %v", output)
		}
		if !strings.Contains(output, "timeout=30") {
			t.Errorf("フィールドtimeoutが含まれていません: %v", output)
		}
	})
}

func TestLogger_Debug(t *testing.T) {
	t.Run("正常系: 基本的なデバッグログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "デバッグ情報"

		// Act
		logger.Debug(ctx, message)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[DEBUG]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
	})

	t.Run("正常系: 複数フィールド付きデバッグログ", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "処理詳細"
		fields := []types.Field{
			{Key: "step", Value: 1},
			{Key: "processing_time", Value: "250ms"},
			{Key: "items_count", Value: 42},
		}

		// Act
		logger.Debug(ctx, message, fields...)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[DEBUG]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
		if !strings.Contains(output, "step=1") {
			t.Errorf("フィールドstepが含まれていません: %v", output)
		}
		if !strings.Contains(output, "processing_time=250ms") {
			t.Errorf("フィールドprocessing_timeが含まれていません: %v", output)
		}
		if !strings.Contains(output, "items_count=42") {
			t.Errorf("フィールドitems_countが含まれていません: %v", output)
		}
	})
}

func TestLogger_FieldFormatting(t *testing.T) {
	t.Run("正常系: 様々な型のフィールド値", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "型テスト"
		fields := []types.Field{
			{Key: "string_val", Value: "テキスト"},
			{Key: "int_val", Value: 123},
			{Key: "bool_val", Value: true},
			{Key: "float_val", Value: 3.14},
		}

		// Act
		logger.Info(ctx, message, fields...)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "string_val=テキスト") {
			t.Errorf("文字列フィールドが正しくフォーマットされていません: %v", output)
		}
		if !strings.Contains(output, "int_val=123") {
			t.Errorf("整数フィールドが正しくフォーマットされていません: %v", output)
		}
		if !strings.Contains(output, "bool_val=true") {
			t.Errorf("真偽値フィールドが正しくフォーマットされていません: %v", output)
		}
		if !strings.Contains(output, "float_val=3.14") {
			t.Errorf("浮動小数点フィールドが正しくフォーマットされていません: %v", output)
		}
	})

	t.Run("正常系: 空のフィールド配列", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		logger := &logger{Logger: createTestLogger(&buf)}
		ctx := context.Background()
		message := "フィールドなし"

		// Act
		logger.Info(ctx, message)

		// Assert
		output := buf.String()
		if !strings.Contains(output, "[INFO]") {
			t.Errorf("ログレベルが含まれていません: %v", output)
		}
		if !strings.Contains(output, message) {
			t.Errorf("メッセージが含まれていません: %v", output)
		}
		// フィールド区切り文字「|」が含まれていないことを確認
		if strings.Contains(output, " |") {
			t.Errorf("フィールド区切り文字が不要に含まれています: %v", output)
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("正常系: ロガーインスタンスの作成", func(t *testing.T) {
		// Act
		logger := New()

		// Assert
		if logger == nil {
			t.Error("ロガーインスタンスがnilです")
		}

		// インターフェースを実装していることを確認
		var _ types.Logger = logger
	})
}

// テスト用のロガーを作成するヘルパー関数
func createTestLogger(buf *bytes.Buffer) *log.Logger {
	return log.New(buf, "[SurveyBot] ", 0) // タイムスタンプなしでテストしやすくする
}