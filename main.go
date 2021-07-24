package main

import (
    "fmt"
    "time"
    "github.com/bwmarrin/discordgo"
    "log"
    "strings"
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
    vcsession *discordgo.VoiceConnection
    HelloWorld = "!helloworld"
    ChannelVoiceJoin = "!vcjoin"
    ChannelVoiceLeave = "!vcleave"
)

func main() {
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

    dg.AddHandler(onMessageCreate) //全てのWSAPIイベントが発生した時のイベントハンドラを追加
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

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    if err != nil {
        log.Println("Error getting channel: ", err)
        return
    }
        fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

    switch {
        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", getenv("YOUR_CLIENT_ID", ""), HelloWorld))://Bot宛に!helloworld コマンドが実行された時
            sendMessage(s, m.ChannelID, "Hello world！")

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", getenv("YOUR_CLIENT_ID", ""), ChannelVoiceJoin)):

            //今いるサーバーのチャンネル情報の一覧を喋らせる処理を書いておきますね
                    //c, err := s.State.Channel(m.ChannelID) //チャンネル取得
            //guildChannels, _ := s.GuildChannels(c.GuildID)
            //var sendText string
            //for _, a := range guildChannels{
                //sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
            //}
            //sendMessage(s, c, sendText) チャンネルの名前、ID、タイプ(通話orテキスト)をBOTが話す

            //VOICE CHANNEL IDには、botを参加させたい通話チャンネルのIDを代入してください
            //コメントアウトされた上記の処理を使うことでチャンネルIDを確認できます
            vcsession, _ = s.ChannelVoiceJoin(c.GuildID, "VOICE_CHANNEL_ID", false, false)
            vcsession.AddHandler(onVoiceReceived) //音声受信時のイベントハンドラ

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", getenv("YOUR_CLIENT_ID", ""), ChannelVoiceLeave)):
            vcsession.Disconnect() //今いる通話チャンネルから抜ける
    }
}

//メッセージを受信した時の、声の初めと終わりにPrintされるようだ
func onVoiceReceived(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
    log.Print("しゃべったあああああ")
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, channelID string, msg string) {
    _, err := s.ChannelMessageSend(c.ID, msg)

    log.Println(">>> " + msg)
    if err != nil {
        log.Println("Error sending message: ", err)
    }
}