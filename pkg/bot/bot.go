package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/Logta/SurveyBot/types"
)

type bot struct {
	session  *discordgo.Session
	handlers []types.Handler
	logger   types.Logger
	config   *types.Config
}

// New creates a new bot instance
func New(config *types.Config, logger types.Logger) (types.Bot, error) {
	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	return &bot{
		session: session,
		logger:  logger,
		config:  config,
	}, nil
}

func (b *bot) RegisterHandler(handler types.Handler) {
	b.handlers = append(b.handlers, handler)
}

func (b *bot) Start(ctx context.Context) error {
	b.session.AddHandler(b.messageCreateHandler)

	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}

	b.logger.Info(ctx, "Bot started successfully")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		b.logger.Info(ctx, "Received interrupt signal, shutting down")
	case <-ctx.Done():
		b.logger.Info(ctx, "Context cancelled, shutting down")
	}

	return b.Stop(ctx)
}

func (b *bot) Stop(ctx context.Context) error {
	if err := b.session.Close(); err != nil {
		b.logger.Error(ctx, "Failed to close Discord session", err)
		return err
	}

	b.logger.Info(ctx, "Bot stopped successfully")
	return nil
}

func (b *bot) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()

	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	b.logger.Debug(ctx, "Received message",
		types.Field{Key: "user", Value: m.Author.Username},
		types.Field{Key: "content", Value: m.Content},
		types.Field{Key: "guild", Value: m.GuildID},
	)

	// Try each handler
	for _, handler := range b.handlers {
		if handler.CanHandle(m.Content) {
			if err := handler.Handle(ctx, s, m); err != nil {
				b.logger.Error(ctx, "Handler failed",
					err,
					types.Field{Key: "handler", Value: handler.Name()},
					types.Field{Key: "message", Value: m.Content},
				)
			}
			return
		}
	}
}