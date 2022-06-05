package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/Logta/SurveyBot/commands"
)

var (
	Help     = "!help"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var (
	stopBot = make(chan bool)
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
	// if m.Content == ServerName {
	// 	g, err := s.Guild(m.GuildID)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(g.Name)
	// 	s.ChannelMessageSend(m.ChannelID, g.Name)
	// }

	commands.SurveyCommands(s, m)
	commands.ShuffleCommands(s, m)
	commands.CouplingCommands(s, m)

	if m.Content == Help {
		commands.SendHelp(s, m)
	}
}
