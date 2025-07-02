package config

import (
	"fmt"
	"os"

	"github.com/Logta/SurveyBot/types"
	"github.com/joho/godotenv"
)

// Load loads configuration from environment variables
func Load() (*types.Config, error) {
	if err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))); err != nil {
		// Ignore error if .env file doesn't exist
	}

	config := &types.Config{
		DiscordToken: getEnv("DISCORD_TOKEN", ""),
		GoEnv:        getEnv("GO_ENV", "development"),
	}

	if config.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
