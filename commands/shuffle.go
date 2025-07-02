package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Logta/SurveyBot/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	Shuffle = "!shuffle"
)

func ShuffleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, Shuffle) {

		temp := regexp.MustCompile(reg).Split(m.Content, -1)
		if len(temp) <= 1 {
			s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでシャッフル項目を記入してください")
			return
		} else if len(temp) <= 2 {
			s.ChannelMessageSend(m.ChannelID, "シャッフル項目は2つ以上記入してください")
			return
		}

		description := ""
		contents := temp[1:]

		contents = utils.FisherYatesShuffle(contents)

		for index, value := range contents {

			e, err := utils.FindEmoji(index)
			if err != nil {
				fmt.Println(err)
				return
			}

			description = description + e + " : " + value + "\n"
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "シャッフル結果",
			Description: description,
			Color:       0x141DB8,
		})
	}
}
