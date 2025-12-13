package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shelldock/shelldock/internal/repo"
)

// Run starts the TUI application
func Run() error {
	manager, err := repo.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize repository manager: %w", err)
	}

	model := NewModel(manager)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

