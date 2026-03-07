package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shelldock/shelldock/internal/repo"
)

func Run() error {
	manager, err := repo.NewManager()
	if err != nil {
		return fmt.Errorf("failed to initialize repository manager: %w", err)
	}

	m := newModel(manager)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
