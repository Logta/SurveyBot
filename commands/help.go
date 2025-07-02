package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SendHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	description := "基本コマンドを上から順に実行することでアンケートが作成できる" + "\n" + "回答項目ごとにスタンプが作成されるため、回答の際には回答項目に対応するスタンプを押下する"

	base_command := ""
	base_command = base_command + Survey + " : " + "アンケート作成を開始する" + "\n"
	base_command = base_command + Title + " : " + "アンケートのタイトルを入力する[改行区切りで入力する]" + "\n"
	base_command = base_command + Content + " : " + "アンケートの回答項目を入力する[改行区切りで入力する]" + "\n"

	confirmation_command := ""
	confirmation_command = confirmation_command + CheckTitle + " : " + "アンケートのタイトルを確認する" + "\n"
	confirmation_command = confirmation_command + CheckState + " : " + "アンケートの設定状況を確認する" + "\n"
	// confirmation_command = confirmation_command + ServerName + " : " + "サーバー名を確認する" + "\n"

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "アンケート機能使い方",
		Description: description,
		Color:       0xA4B814,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "基本コマンド", Value: base_command, Inline: true},
			&discordgo.MessageEmbedField{Name: "キャンセルコマンド", Value: Cancel + " : " + "アンケートの作成を中止する" + "\n", Inline: true},
			&discordgo.MessageEmbedField{Name: "確認コマンド", Value: confirmation_command, Inline: false},
		},
	})

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "シャッフル機能使い方",
		Description: Shuffle + " : " + "与えられた項目をシャッフルする[項目は改行区切りで入力する]" + "\n",
		Color:       0xA4B814,
	})

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "カップリング機能使い方",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "基本コマンド", Value: Coupling + " : " + "与えられた項目で組み合わせを作る。組み合わせる集合は改行で区切り、集合内はカンマ区切りで入力。" + "\n", Inline: true},
		},
		Description: "基本コマンドに改行で区切った、カンマ区切りの集合を指定することで集合同士の要素の組み合わせを作成します。",
		Color:       0xA4B814,
	})
}
