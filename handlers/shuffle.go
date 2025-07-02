package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Logta/SurveyBot/types"
	"github.com/bwmarrin/discordgo"
)

type shuffleHandler struct {
	shuffler      types.Shuffler
	emojiProvider types.EmojiProvider
	logger        types.Logger
	regexPattern  *regexp.Regexp
}

// NewShuffleHandler creates a new shuffle command handler
func NewShuffleHandler(shuffler types.Shuffler, emojiProvider types.EmojiProvider, logger types.Logger) types.Handler {
	return &shuffleHandler{
		shuffler:      shuffler,
		emojiProvider: emojiProvider,
		logger:        logger,
		regexPattern:  regexp.MustCompile(`\r\n|\n| |,`),
	}
}

func (h *shuffleHandler) Name() string {
	return "ShuffleHandler"
}

func (h *shuffleHandler) CanHandle(command string) bool {
	return strings.HasPrefix(command, string(types.CmdShuffle))
}

func (h *shuffleHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error {
	parts := h.regexPattern.Split(m.Content, -1)
	if len(parts) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでシャッフル項目を記入してください")
		return err
	}

	if len(parts) <= 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "シャッフル項目は2つ以上記入してください")
		return err
	}

	items := parts[1:]
	shuffledItems := h.shuffler.Shuffle(ctx, items)

	return h.createShuffleEmbed(ctx, s, m, shuffledItems)
}

func (h *shuffleHandler) createShuffleEmbed(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, items []string) error {
	description := ""
	maxEmojis := h.emojiProvider.GetMaxEmojis()

	if len(items) > maxEmojis {
		return fmt.Errorf("too many items: %d (max: %d)", len(items), maxEmojis)
	}

	for i, item := range items {
		emoji, err := h.emojiProvider.GetEmoji(ctx, i)
		if err != nil {
			h.logger.Error(ctx, "Failed to get emoji", err)
			return err
		}

		description += fmt.Sprintf("%s : %s\n", emoji, item)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "シャッフル結果",
		Description: description,
		Color:       0x141DB8,
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}
