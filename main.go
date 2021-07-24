package main

import (
    "fmt"
    "log"
    "os"
    "errors"

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

var(
    stopBot = make(chan bool)
    ServerName = "!servername"
    ChannelVoiceJoin = "!vcjoin"
    ChannelVoiceLeave = "!vcleave"
	state = false
	title = ""
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
	if m.Content == "!survey" {
		state = true
		s.ChannelMessageSend(m.ChannelID, "アンケートのタイトルを入力してください")
	}

	//途中で止める用のコマンド
	if m.Content == "!cancel" {
		state = false
		s.ChannelMessageSend(m.ChannelID, "アンケート作成をキャンセルしました")
	}

	//アンケート情報入力スタートのコマンド
	if m.Content == "!title" {
		state = true
		s.ChannelMessageSend(m.ChannelID, "アンケートのタイトルを入力してください")
	}

	if m.Content == "!state check" {
		if state && title != "" {
			s.ChannelMessageSend(m.ChannelID, "アンケート内容を記入してください")
		} else if state && title == "" {
			s.ChannelMessageSend(m.ChannelID, "アンケートタイトルを入力してください")
		} else {
			s.ChannelMessageSend(m.ChannelID, "アンケートは開始されていません")
		}
	}

	if m.Content =="test" {

		mes, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "あんけ",
			Description: "Don't ever talk to me or my son ever again.",
		})

		emoji, err := FindEmoji(1)
		if err != nil{
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

func FindEmoji(num int) (string, error){

	switch(num){
		case 0 : return "0️⃣", nil
		case 1 : return "1️⃣", nil
		case 2 : return "2️⃣", nil
		case 3 : return "3️⃣", nil
		case 4 : return "4️⃣", nil
		case 5 : return "5️⃣", nil
		case 6 : return "6️⃣", nil
		case 7 : return "7️⃣", nil
		case 8 : return "8️⃣", nil
		case 9 : return "9️⃣", nil
		case 10 : return "🔟", nil
		default: return "", errors.New("絵文字が見つかりません")
	}
}