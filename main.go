package main

import (
	"context"
	"log"

	"github.com/Logta/SurveyBot/pkg/bot"
	"github.com/Logta/SurveyBot/pkg/config"
	"github.com/Logta/SurveyBot/pkg/logger"
	"github.com/Logta/SurveyBot/pkg/state"
	"github.com/Logta/SurveyBot/handlers"
	"github.com/Logta/SurveyBot/utils"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.New()

	// Initialize dependencies
	stateManager := state.NewMemoryStateManager()
	emojiProvider := utils.NewEmojiProvider()
	shuffler := utils.NewShuffler()
	coupler := utils.NewCoupler()

	// Create bot instance
	b, err := bot.New(cfg, logger)
	if err != nil {
		logger.Error(ctx, "Failed to create bot", err)
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Register handlers
	b.RegisterHandler(handlers.NewSurveyHandler(stateManager, emojiProvider, logger))
	b.RegisterHandler(handlers.NewShuffleHandler(shuffler, emojiProvider, logger))
	b.RegisterHandler(handlers.NewCouplingHandler(coupler, emojiProvider, logger))
	b.RegisterHandler(handlers.NewHelpHandler(logger))

	// Start bot
	if err := b.Start(ctx); err != nil {
		logger.Error(ctx, "Bot failed to start", err)
		log.Fatalf("Bot failed to start: %v", err)
	}
}
