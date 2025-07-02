package types

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

// Config holds application configuration
type Config struct {
	DiscordToken string
	GoEnv        string
}

// SurveyState represents the state of a survey creation
type SurveyState struct {
	Active bool
	Title  string
}

// Command represents a Discord command
type Command string

const (
	CmdHelp       Command = "!help"
	CmdSurvey     Command = "!survey"
	CmdTitle      Command = "!title"
	CmdContent    Command = "!content"
	CmdCancel     Command = "!cancel"
	CmdCheckState Command = "!check state"
	CmdCheckTitle Command = "!check title"
	CmdShuffle    Command = "!shuffle"
	CmdCoupling   Command = "!coupling"
)

// Handler defines the interface for command handlers
type Handler interface {
	Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) error
	CanHandle(command string) bool
	Name() string
}

// StateManager manages survey state
type StateManager interface {
	GetState(ctx context.Context, guildID string) (*SurveyState, error)
	SetState(ctx context.Context, guildID string, state *SurveyState) error
	ClearState(ctx context.Context, guildID string) error
}

// EmojiProvider provides emoji utilities
type EmojiProvider interface {
	GetEmoji(ctx context.Context, index int) (string, error)
	GetMaxEmojis() int
}

// Shuffler provides shuffle functionality
type Shuffler interface {
	Shuffle(ctx context.Context, items []string) []string
}

// Coupler provides coupling functionality
type Coupler interface {
	Couple(ctx context.Context, itemSets [][]string) ([][]string, error)
}

// Logger defines logging interface
type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)
}

// Field represents a structured log field
type Field struct {
	Key   string
	Value interface{}
}

// Bot represents the main bot instance
type Bot interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RegisterHandler(handler Handler)
}
