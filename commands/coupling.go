package commands

import (
	"regexp"
	"strings"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/Logta/SurveyBot/utils"
)

var (
	Coupling = "!coupling"
)

func CouplingCommands(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, Coupling) {

		temp := regexp.MustCompile(regIndention).Split(m.Content, -1)
		if len(temp) <= 2 {
			s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでカップリング対象の集合を2つ以上を記入してください")
			return
		}

		description := ""
		lines := temp[1:]
		itemSets := utils.GetItemSets(temp, regCSV)

		var base [][]string
		result := utils.Coupling(itemSets, base)

		for index, value := range result {

			e, err := utils.FindEmoji(index)
			if err != nil {
				fmt.Println(err)
				return
			}

			description = description + e + " " + value[0] + " : " + 
			strings.Join(value[1:], ",") + "\n"
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "シャッフル結果",
			Description: description,
			Color:       0x141DB8,
		})
	}
}