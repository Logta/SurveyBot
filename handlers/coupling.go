package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/Logta/SurveyBot/types"
	"github.com/Logta/SurveyBot/utils"
)

type couplingHandler struct {
	coupler       types.Coupler
	emojiProvider types.EmojiProvider
	logger        types.Logger
	lineRegex     *regexp.Regexp
}

// NewCouplingHandler creates a new coupling command handler
func NewCouplingHandler(coupler types.Coupler, emojiProvider types.EmojiProvider, logger types.Logger) types.Handler {
	return &couplingHandler{
		coupler:       coupler,
		emojiProvider: emojiProvider,
		logger:        logger,
		lineRegex:     regexp.MustCompile(`\r\n|\n`),
	}
}

func (h *couplingHandler) Name() string {
	return "CouplingHandler"
}

func (h *couplingHandler) CanHandle(command string) bool {
	return strings.HasPrefix(command, string(types.CmdCoupling))
}

func (h *couplingHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error {
	lines := h.lineRegex.Split(m.Content, -1)
	if len(lines) <= 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでカップリング対象の集合を2つ以上を記入してください")
		return err
	}

	h.logger.Debug(ctx, "Processing coupling command",
		types.Field{Key: "lines_count", Value: len(lines)},
		types.Field{Key: "content", Value: strings.Join(lines[1:], ",")},
	)

	itemSets := utils.ParseItemSets(lines[1:], ",")
	h.logger.Debug(ctx, "Parsed item sets", types.Field{Key: "sets_count", Value: len(itemSets)})

	result, err := h.coupler.Couple(ctx, itemSets)
	if err != nil {
		h.logger.Error(ctx, "Failed to couple items", err)
		return err
	}

	h.logger.Debug(ctx, "Coupling completed", types.Field{Key: "result_count", Value: len(result)})

	return h.createCouplingEmbed(ctx, s, m, result)
}

func (h *couplingHandler) createCouplingEmbed(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, couples [][]string) error {
	description := ""
	maxEmojis := h.emojiProvider.GetMaxEmojis()

	if len(couples) > maxEmojis {
		return fmt.Errorf("too many couples: %d (max: %d)", len(couples), maxEmojis)
	}

	for i, couple := range couples {
		emoji, err := h.emojiProvider.GetEmoji(ctx, i)
		if err != nil {
			h.logger.Error(ctx, "Failed to get emoji", err)
			return err
		}

		if len(couple) == 0 {
			continue
		}

		// Format: "emoji leader : member1,member2,member3"
		leader := couple[0]
		members := strings.Join(couple[1:], ",")
		description += fmt.Sprintf("%s %s : %s\n", emoji, leader, members)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "カップリング結果",
		Description: description,
		Color:       0x141DB8,
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	return err
}