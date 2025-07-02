package main

import (
	"context"
	"testing"

	"github.com/Logta/SurveyBot/types"
	"github.com/Logta/SurveyBot/pkg/state"
	"github.com/Logta/SurveyBot/utils"
	"github.com/Logta/SurveyBot/handlers"
	"github.com/Logta/SurveyBot/pkg/logger"
)

// TestHelper provides common test utilities and mocks
type TestHelper struct {
	StateManager  types.StateManager
	EmojiProvider types.EmojiProvider
	Shuffler      types.Shuffler
	Coupler       types.Coupler
	Logger        types.Logger
}

// NewTestHelper creates a new test helper with real implementations
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		StateManager:  state.NewMemoryStateManager(),
		EmojiProvider: utils.NewEmojiProvider(),
		Shuffler:      utils.NewShuffler(),
		Coupler:       utils.NewCoupler(),
		Logger:        logger.New(),
	}
}

// CreateSurveyHandler creates a survey handler for testing
func (h *TestHelper) CreateSurveyHandler() types.Handler {
	return handlers.NewSurveyHandler(h.StateManager, h.EmojiProvider, h.Logger)
}

// CreateShuffleHandler creates a shuffle handler for testing
func (h *TestHelper) CreateShuffleHandler() types.Handler {
	return handlers.NewShuffleHandler(h.Shuffler, h.EmojiProvider, h.Logger)
}

// CreateCouplingHandler creates a coupling handler for testing
func (h *TestHelper) CreateCouplingHandler() types.Handler {
	return handlers.NewCouplingHandler(h.Coupler, h.EmojiProvider, h.Logger)
}

// CreateHelpHandler creates a help handler for testing
func (h *TestHelper) CreateHelpHandler() types.Handler {
	return handlers.NewHelpHandler(h.Logger)
}

// SetupSurveyState sets up a survey state for testing
func (h *TestHelper) SetupSurveyState(t *testing.T, guildID string, active bool, title string) {
	ctx := context.Background()
	state := &types.SurveyState{
		Active: active,
		Title:  title,
	}
	
	err := h.StateManager.SetState(ctx, guildID, state)
	if err != nil {
		t.Fatalf("Failed to setup survey state: %v", err)
	}
}

// VerifySurveyState verifies the survey state matches expectations
func (h *TestHelper) VerifySurveyState(t *testing.T, guildID string, expectedActive bool, expectedTitle string) {
	ctx := context.Background()
	state, err := h.StateManager.GetState(ctx, guildID)
	if err != nil {
		t.Fatalf("Failed to get survey state: %v", err)
	}
	
	if state.Active != expectedActive {
		t.Errorf("Survey active state mismatch: got %v, want %v", state.Active, expectedActive)
	}
	
	if state.Title != expectedTitle {
		t.Errorf("Survey title mismatch: got %v, want %v", state.Title, expectedTitle)
	}
}