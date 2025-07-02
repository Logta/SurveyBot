package handlers

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/Logta/SurveyBot/types"
)

type helpHandler struct {
	logger types.Logger
}

// NewHelpHandler creates a new help command handler
func NewHelpHandler(logger types.Logger) types.Handler {
	return &helpHandler{
		logger: logger,
	}
}

func (h *helpHandler) Name() string {
	return "HelpHandler"
}

func (h *helpHandler) CanHandle(command string) bool {
	return command == string(types.CmdHelp)
}

func (h *helpHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return h.sendHelpEmbeds(ctx, s, m)
}

func (h *helpHandler) sendHelpEmbeds(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error {
	// Survey help embed
	surveyDescription := "基本コマンドを上から順に実行することでアンケートが作成できる" + "\n" + "回答項目ごとにスタンプが作成されるため、回答の際には回答項目に対応するスタンプを押下する"

	baseCommands := ""
	baseCommands += string(types.CmdSurvey) + " : " + "アンケート作成を開始する" + "\n"
	baseCommands += string(types.CmdTitle) + " : " + "アンケートのタイトルを入力する[改行区切りで入力する]" + "\n"
	baseCommands += string(types.CmdContent) + " : " + "アンケートの回答項目を入力する[改行区切りで入力する]" + "\n"

	confirmationCommands := ""
	confirmationCommands += string(types.CmdCheckTitle) + " : " + "アンケートのタイトルを確認する" + "\n"
	confirmationCommands += string(types.CmdCheckState) + " : " + "アンケートの設定状況を確認する" + "\n"

	surveyEmbed := &discordgo.MessageEmbed{
		Title:       "アンケート機能使い方",
		Description: surveyDescription,
		Color:       0xA4B814,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "基本コマンド", Value: baseCommands, Inline: true},
			{Name: "キャンセルコマンド", Value: string(types.CmdCancel) + " : " + "アンケートの作成を中止する" + "\n", Inline: true},
			{Name: "確認コマンド", Value: confirmationCommands, Inline: false},
		},
	}

	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, surveyEmbed); err != nil {
		return err
	}

	// Shuffle help embed
	shuffleEmbed := &discordgo.MessageEmbed{
		Title:       "シャッフル機能使い方",
		Description: string(types.CmdShuffle) + " : " + "与えられた項目をシャッフルする[項目は改行区切りで入力する]" + "\n",
		Color:       0xA4B814,
	}

	if _, err := s.ChannelMessageSendEmbed(m.ChannelID, shuffleEmbed); err != nil {
		return err
	}

	// Coupling help embed
	couplingEmbed := &discordgo.MessageEmbed{
		Title:       "カップリング機能使い方",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "基本コマンド", Value: string(types.CmdCoupling) + " : " + "与えられた項目で組み合わせを作る。組み合わせる集合は改行で区切り、集合内はカンマ区切りで入力。" + "\n", Inline: true},
		},
		Description: "基本コマンドに改行で区切った、カンマ区切りの集合を指定することで集合同士の要素の組み合わせを作成します。",
		Color:       0xA4B814,
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, couplingEmbed)
	return err
}