package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Logta/SurveyBot/types"
	"github.com/bwmarrin/discordgo"
)

type surveyHandler struct {
	stateManager  types.StateManager
	emojiProvider types.EmojiProvider
	logger        types.Logger
	regexPattern  *regexp.Regexp
}

// NewSurveyHandler creates a new survey command handler
func NewSurveyHandler(stateManager types.StateManager, emojiProvider types.EmojiProvider, logger types.Logger) types.Handler {
	return &surveyHandler{
		stateManager:  stateManager,
		emojiProvider: emojiProvider,
		logger:        logger,
		regexPattern:  regexp.MustCompile(`\r\n|\n| |,`),
	}
}

func (h *surveyHandler) Name() string {
	return "SurveyHandler"
}

func (h *surveyHandler) CanHandle(command string) bool {
	return strings.HasPrefix(command, string(types.CmdSurvey)) ||
		strings.HasPrefix(command, string(types.CmdTitle)) ||
		strings.HasPrefix(command, string(types.CmdContent)) ||
		command == string(types.CmdCancel) ||
		command == string(types.CmdCheckState) ||
		command == string(types.CmdCheckTitle)
}

func (h *surveyHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error {
	guildID := m.GuildID
	if guildID == "" {
		guildID = "dm" // Handle direct messages
	}

	switch {
	case m.Content == string(types.CmdSurvey):
		return h.handleSurveyStart(ctx, s, m, guildID)

	case m.Content == string(types.CmdCancel):
		return h.handleCancel(ctx, s, m, guildID)

	case m.Content == string(types.CmdCheckState):
		return h.handleCheckState(ctx, s, m, guildID)

	case m.Content == string(types.CmdCheckTitle):
		return h.handleCheckTitle(ctx, s, m, guildID)

	case strings.HasPrefix(m.Content, string(types.CmdTitle)):
		return h.handleTitle(ctx, s, m, guildID)

	case strings.HasPrefix(m.Content, string(types.CmdContent)):
		return h.handleContent(ctx, s, m, guildID)
	}

	return nil
}

func (h *surveyHandler) handleSurveyStart(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	state := &types.SurveyState{
		Active: true,
		Title:  "",
	}

	if err := h.stateManager.SetState(ctx, guildID, state); err != nil {
		h.logger.Error(ctx, "Failed to set survey state", err)
		return err
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "アンケートのタイトルを入力してください")
	return err
}

func (h *surveyHandler) handleCancel(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	if err := h.stateManager.ClearState(ctx, guildID); err != nil {
		h.logger.Error(ctx, "Failed to clear survey state", err)
		return err
	}

	_, err := s.ChannelMessageSend(m.ChannelID, "アンケート作成をキャンセルしました")
	return err
}

func (h *surveyHandler) handleCheckState(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	state, err := h.stateManager.GetState(ctx, guildID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get survey state", err)
		return err
	}

	var message string
	if state.Active && state.Title != "" {
		message = "アンケート内容を記入してください"
	} else if state.Active && state.Title == "" {
		message = "アンケートタイトルを入力してください"
	} else {
		message = "アンケートは開始されていません"
	}

	_, err = s.ChannelMessageSend(m.ChannelID, message)
	return err
}

func (h *surveyHandler) handleCheckTitle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	state, err := h.stateManager.GetState(ctx, guildID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get survey state", err)
		return err
	}

	message := "現在設定されてるタイトル\n" + state.Title
	_, err = s.ChannelMessageSend(m.ChannelID, message)
	return err
}

func (h *surveyHandler) handleTitle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	state, err := h.stateManager.GetState(ctx, guildID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get survey state", err)
		return err
	}

	if !state.Active {
		return nil // Ignore if survey is not active
	}

	parts := h.regexPattern.Split(m.Content, -1)
	if len(parts) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んでタイトルを記入してください")
		return err
	}

	state.Title = parts[1]
	if err := h.stateManager.SetState(ctx, guildID, state); err != nil {
		h.logger.Error(ctx, "Failed to update survey state", err)
		return err
	}

	return nil
}

func (h *surveyHandler) handleContent(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, guildID string) error {
	state, err := h.stateManager.GetState(ctx, guildID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get survey state", err)
		return err
	}

	if !state.Active {
		return nil // Ignore if survey is not active
	}

	parts := h.regexPattern.Split(m.Content, -1)
	if len(parts) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "コマンドの後に改行を挟んで回答項目を記入してください")
		return err
	}

	return h.createSurveyEmbed(ctx, s, m, state.Title, parts[1:])
}

func (h *surveyHandler) createSurveyEmbed(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, title string, options []string) error {
	description := ""
	indices := []int{}

	maxEmojis := h.emojiProvider.GetMaxEmojis()
	if len(options) > maxEmojis {
		return fmt.Errorf("too many options: %d (max: %d)", len(options), maxEmojis)
	}

	for i, option := range options {
		emoji, err := h.emojiProvider.GetEmoji(ctx, i+1) // Start from 1
		if err != nil {
			h.logger.Error(ctx, "Failed to get emoji", err)
			return err
		}

		indices = append(indices, i+1)
		description += fmt.Sprintf("%s : %s\n", emoji, option)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x141DB8,
	}

	message, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	// Add reactions
	for _, idx := range indices {
		emoji, err := h.emojiProvider.GetEmoji(ctx, idx)
		if err != nil {
			h.logger.Error(ctx, "Failed to get emoji for reaction", err)
			continue
		}

		if err := s.MessageReactionAdd(m.ChannelID, message.ID, emoji); err != nil {
			h.logger.Error(ctx, "Failed to add reaction", err)
		}
	}

	return nil
}
