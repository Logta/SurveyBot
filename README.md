# SurveyBot

Discord 上でアンケート、シャッフル、チーム編成を行うための Bot

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 機能

### アンケート作成

- 絵文字リアクション付きの投票作成
- 複数選択肢のアンケート
- インタラクティブな回答収集

### シャッフル・ランダム化

- アイテムリストのランダム並び替え
- 簡単なチーム編成

### チーム編成

- 複数グループからのチームペアリング
- グループ間メンバー配分

## クイックスタート

### Bot 招待

[SurveyBot をサーバーに招待](https://discord.com/oauth2/authorize?client_id=868454195953561610&scope=bot&permissions=0)

### 基本コマンド

```
!help          # 利用可能なコマンドを表示
!survey        # アンケート作成を開始
!shuffle       # アイテムリストをシャッフル
!coupling      # チーム編成を実行
```

## 使用例

### アンケート作成

```
!survey
!title
好きなプログラミング言語は？

!content
Go
TypeScript
Rust
Python
```

### シャッフル

```
!shuffle
田中
佐藤
鈴木
高橋
```

### チーム編成

```
!coupling
フロントエンド,バックエンド,DevOps
田中,佐藤
鈴木,高橋
山田,佐々木
```

## 開発

### 必要な環境

- [mise](https://mise.jdx.dev/) (ツール管理)
- Go 1.24+
- Discord Bot Token

### セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/Logta/SurveyBot.git
cd SurveyBot

# 依存関係のインストール
mise install
mise run deps

# 環境ファイルの初期化
mise run init-env
# .envファイルにDISCORD_TOKENを設定
```

### 開発ワークフロー

```bash
# 開発サイクル
mise run dev           # フォーマット、リント、テスト

# テスト
mise run test          # 全テスト実行
mise run test-coverage # カバレッジレポート生成
mise run test-race     # 競合状態検出

# ビルド
mise run build         # 開発用ビルド
mise run build-release # 本番用最適化ビルド

# 実行
mise run run           # 開発モード
mise run run-prod      # 本番モード
```

## プロジェクト構成

```
SurveyBot/
├── types/              # 型定義とインターフェース
├── pkg/
│   ├── bot/           # Bot実装
│   ├── config/        # 設定管理
│   ├── logger/        # ログ機能
│   └── state/         # 状態管理
├── handlers/          # コマンドハンドラー
├── utils/             # ユーティリティ関数
└── tests/             # テストヘルパーと統合テスト
```

## テスト

t-wada スタイルのテストと AAA パターンを採用:

```go
func TestEmojiProvider_GetEmoji(t *testing.T) {
    t.Run("正常系: 有効なインデックスで絵文字を取得", func(t *testing.T) {
        // Arrange
        provider := NewEmojiProvider()
        ctx := context.Background()

        // Act
        result, err := provider.GetEmoji(ctx, 0)

        // Assert
        if err != nil {
            t.Errorf("期待していないエラーが発生: %v", err)
        }
        if result != "0️⃣" {
            t.Errorf("絵文字が期待値と異なります: got %v, want %v", result, "0️⃣")
        }
    })
}
```

### テストカバレッジ

- **utils**: 100%
- **pkg/logger**: 100%
- **pkg/state**: 100%
- **統合テスト**: 全テスト通過

## デプロイ

### Heroku

```bash
# Herokuにデプロイ
mise run heroku-deploy

# ワーカーのスケール（必要に応じて）
heroku ps:scale worker=1
```

### Docker

```bash
# Dockerイメージのビルド
mise run docker-build

# Dockerで実行
mise run docker-run
```

### 環境変数

```bash
DISCORD_TOKEN=your_discord_bot_token
GO_ENV=production  # オプション、デフォルトはdevelopment
```

### 開発ガイドライン

- Go の規約と`gofmt`フォーマットに従う
- 新機能にはテストを書く
- わかりやすいコミットメッセージを使用
- ユーザー向け変更時はドキュメントを更新

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています - 詳細は[LICENSE](LICENSE)ファイルを参照してください。
