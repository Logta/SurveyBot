package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var (
	stopBot    = make(chan bool)
	ServerName = "!servername"
	Help       = "!help"
	Survey     = "!survey"
	Title      = "!title"
	Content    = "!content"
	Cancel     = "!cancel"
	CheckState = "!check state"
	CheckTitle = "!check title"

	state = false
	title = ""
	reg   = "\r\n|\n| |,"
)

func main() {
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))

	//Discordのセッションを作成
	dg, err := discordgo.New("Bot " + getenv("DISCORD_TOKEN", ""))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	dg.AddHandler(messageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加

	// websocketを開いてlistening開始
	err = dg.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer dg.Close()

	fmt.Println("Listening...")
	<-stopBot //プログラムが終了しないようロック
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Server名を取得して返します。
	if m.Content == ServerName {
		g, err := s.Guild(m.GuildID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(g.Name)
		s.ChannelMessageSend(m.ChannelID, g.Name)
	}

	//アンケート情報入力スタートのコマンド
	if m.Content == Survey {
		state = true
		s.ChannelMessageSend(m.ChannelID, "アンケートのタイトルを入力してください")
	}

	//途中で止める用のコマンド
	if m.Content == Cancel {
		state = false
		title = ""
		s.ChannelMessageSend(m.ChannelID, "アンケート作成をキャンセルしました")
	}

	if m.Content == CheckState {
		if state && title != "" {
			s.ChannelMessageSend(m.ChannelID, "アンケート内容を記入してください")
		} else if state && title == "" {
			s.ChannelMessageSend(m.ChannelID, "アンケートタイトルを入力してください")
		} else {
			s.ChannelMessageSend(m.ChannelID, "アンケートは開始されていません")
		}
	}

	if m.Content == CheckTitle {
		s.ChannelMessageSend(m.ChannelID, "現在設定されてるタイトル\n"+title)
	}

	if strings.HasPrefix(m.Content, Title) {
		if !state {
			return
		}

		temp := regexp.MustCompile(reg).Split(m.Content, -1)
		if len(temp) <= 1 {
			s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでタイトルを記入してください")
		}

		title = temp[1]
	}

	if strings.HasPrefix(m.Content, Content) {
		if !state {
			return
		}

		temp := regexp.MustCompile(reg).Split(m.Content, -1)
		if len(temp) <= 1 {
			s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んで回答項目を記入してください")
		}

		description := ""
		numbers := []int{}

		for index, value := range temp {
			if index == 0 {
				continue
			}

			e, err := FindEmoji(index)
			if err != nil {
				fmt.Println(err)
				return
			}

			numbers = append(numbers, index)
			description = description + e + " : " + value + "\n"
		}

		mes, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       title,
			Description: description,
			Color:       0x141DB8,
		})

		for _, i := range numbers {
			emoji, err := FindEmoji(i)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = s.MessageReactionAdd(m.ChannelID, mes.ID, emoji)
			if err != nil {
				fmt.Println("Error logging in")
				fmt.Println(err)
			}
		}
	}

	if m.Content == Help {
		description := "基本コマンドを上から順に実行することでアンケートが作成できる" + "\n" + "回答項目ごとにスタンプが作成されるため、回答の際には回答項目に対応するスタンプを押下する"

		base_command := ""
		base_command = base_command + Survey + " : " + "アンケート作成を開始する" + "\n"
		base_command = base_command + Title + " : " + "アンケートのタイトルを入力する[改行区切りで入力する]" + "\n"
		base_command = base_command + Content + " : " + "アンケートの回答項目を入力する[改行区切りで入力する]" + "\n"

		confirmation_command := ""
		confirmation_command = confirmation_command + CheckTitle + " : " + "アンケートのタイトルを確認する" + "\n"
		confirmation_command = confirmation_command + CheckState + " : " + "アンケートの設定状況を確認する" + "\n"
		confirmation_command = confirmation_command + ServerName + " : " + "サーバー名を確認する" + "\n"

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "使い方",
			Description: description,
			Color:       0xA4B814,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{Name: "基本コマンド", Value: base_command, Inline: true},
				&discordgo.MessageEmbedField{Name: "キャンセルコマンド", Value: Cancel + " : " + "アンケートの作成を中止する" + "\n", Inline: true},
				&discordgo.MessageEmbedField{Name: "確認コマンド", Value: confirmation_command, Inline: false},
			},
		})
	}
}

func FindEmoji(num int) (string, error) {

	switch num {
	case 0:
		return "0️⃣", nil
	case 1:
		return "1️⃣", nil
	case 2:
		return "2️⃣", nil
	case 3:
		return "3️⃣", nil
	case 4:
		return "4️⃣", nil
	case 5:
		return "5️⃣", nil
	case 6:
		return "6️⃣", nil
	case 7:
		return "7️⃣", nil
	case 8:
		return "8️⃣", nil
	case 9:
		return "9️⃣", nil
	case 10:
		return "🔟", nil
	default:
		return "", errors.New("絵文字が見つかりません")
	}
}
