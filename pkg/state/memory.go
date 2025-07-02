package state

import (
	"context"
	"fmt"
	"sync"

	"github.com/Logta/SurveyBot/types"
)

type memoryStateManager struct {
	mu     sync.RWMutex
	states map[string]*types.SurveyState
}

// NewMemoryStateManager creates a new in-memory state manager
func NewMemoryStateManager() types.StateManager {
	return &memoryStateManager{
		states: make(map[string]*types.SurveyState),
	}
}

func (m *memoryStateManager) GetState(ctx context.Context, guildID string) (*types.SurveyState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if state, exists := m.states[guildID]; exists {
		// Return a copy to avoid race conditions
		return &types.SurveyState{
			Active: state.Active,
			Title:  state.Title,
		}, nil
	}

	return &types.SurveyState{}, nil
}

func (m *memoryStateManager) SetState(ctx context.Context, guildID string, state *types.SurveyState) error {
	if state == nil {
		return fmt.Errorf("state cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.states[guildID] = &types.SurveyState{
		Active: state.Active,
		Title:  state.Title,
	}

	return nil
}

func (m *memoryStateManager) ClearState(ctx context.Context, guildID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.states, guildID)
	return nil
}