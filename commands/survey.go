package commands

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
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

func SurveyCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
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

			e, err := utils.FindEmoji(index)
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
			emoji, err := utils.FindEmoji(i)
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
}
