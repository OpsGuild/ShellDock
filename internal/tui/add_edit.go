package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shelldock/shelldock/internal/repo"
)

type addEditModel struct {
	manager *repo.Manager
	isEdit  bool
	cmdSet  *repo.CommandSet
	step    int // 0: name, 1: description, 2: version, 3: commands
	cmdIdx  int // current command being edited
	field   int // 0: desc, 1: command
	input   string
	err     error
}

func newAddEditModel(manager *repo.Manager, isEdit bool, cmdSet *repo.CommandSet) *addEditModel {
	if cmdSet == nil {
		cmdSet = &repo.CommandSet{
			Name:        "",
			Description: "",
			Version:     "1.0.0",
			Commands:    []repo.Command{},
		}
	}
	return &addEditModel{
		manager: manager,
		isEdit:  isEdit,
		cmdSet:  cmdSet,
		step:    0,
		cmdIdx:  -1,
		field:   0,
		input:   "",
	}
}

func (m *addEditModel) Init() tea.Cmd {
	return nil
}

func (m *addEditModel) Update(msg tea.Msg) (*addEditModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return nil, tea.Quit
		case "enter":
			return m.handleEnter()
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "tab":
			if m.step == 3 {
				m.field = 1 - m.field // toggle
			}
		case "ctrl+s":
			return m.save()
		default:
			if msg.Type == tea.KeyRunes {
				m.input += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m *addEditModel) handleEnter() (*addEditModel, tea.Cmd) {
	switch m.step {
	case 0: // name
		if strings.TrimSpace(m.input) != "" {
			m.cmdSet.Name = strings.TrimSpace(m.input)
			m.input = m.cmdSet.Description
			m.step = 1
		}
	case 1: // description
		m.cmdSet.Description = strings.TrimSpace(m.input)
		m.input = m.cmdSet.Version
		m.step = 2
	case 2: // version
		m.cmdSet.Version = strings.TrimSpace(m.input)
		if m.cmdSet.Version == "" {
			m.cmdSet.Version = "1.0.0"
		}
		m.step = 3
		m.cmdIdx = 0
		m.input = ""
		if len(m.cmdSet.Commands) == 0 {
			m.cmdSet.Commands = append(m.cmdSet.Commands, repo.Command{})
		}
	case 3: // commands
		if m.field == 0 { // description
			if strings.TrimSpace(m.input) != "" {
				if m.cmdIdx >= len(m.cmdSet.Commands) {
					m.cmdSet.Commands = append(m.cmdSet.Commands, repo.Command{})
				}
				m.cmdSet.Commands[m.cmdIdx].Description = strings.TrimSpace(m.input)
				m.input = m.cmdSet.Commands[m.cmdIdx].Command
				m.field = 1
			}
		} else { // command
			if strings.TrimSpace(m.input) != "" {
				if m.cmdIdx >= len(m.cmdSet.Commands) {
					m.cmdSet.Commands = append(m.cmdSet.Commands, repo.Command{})
				}
				m.cmdSet.Commands[m.cmdIdx].Command = strings.TrimSpace(m.input)
				m.cmdIdx++
				m.input = ""
				m.field = 0
				if m.cmdIdx >= len(m.cmdSet.Commands) {
					m.cmdSet.Commands = append(m.cmdSet.Commands, repo.Command{})
				}
			}
		}
	}
	return m, nil
}

func (m *addEditModel) save() (*addEditModel, tea.Cmd) {
	if m.cmdSet.Name == "" {
		m.err = fmt.Errorf("name is required")
		return m, nil
	}
	
	// Remove empty commands
	var validCommands []repo.Command
	for _, cmd := range m.cmdSet.Commands {
		if cmd.Description != "" && cmd.Command != "" {
			validCommands = append(validCommands, cmd)
		}
	}
	m.cmdSet.Commands = validCommands

	if err := m.manager.GetLocalRepo().SaveCommandSet(m.cmdSet, ""); err != nil {
		m.err = err
		return m, nil
	}

	return nil, tea.Quit
}

func (m *addEditModel) View() string {
	var b strings.Builder
	title := "Add Command Set"
	if m.isEdit {
		title = "Edit Command Set"
	}
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s ", title)))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(fmt.Sprintf("Error: %v\n\n", m.err)))
	}

	switch m.step {
	case 0:
		b.WriteString("Name: " + m.input + "█\n")
		b.WriteString("\nPress Enter to continue")
	case 1:
		b.WriteString("Name: " + m.cmdSet.Name + "\n")
		b.WriteString("Description: " + m.input + "█\n")
		b.WriteString("\nPress Enter to continue")
	case 2:
		b.WriteString("Name: " + m.cmdSet.Name + "\n")
		b.WriteString("Description: " + m.cmdSet.Description + "\n")
		b.WriteString("Version: " + m.input + "█\n")
		b.WriteString("\nPress Enter to add commands")
	case 3:
		b.WriteString("Name: " + m.cmdSet.Name + "\n")
		b.WriteString("Description: " + m.cmdSet.Description + "\n")
		b.WriteString("Version: " + m.cmdSet.Version + "\n\n")
		b.WriteString("Commands:\n")
		
		for i := 0; i <= m.cmdIdx && i < len(m.cmdSet.Commands); i++ {
			cmd := m.cmdSet.Commands[i]
			if i == m.cmdIdx {
				if m.field == 0 {
					b.WriteString(fmt.Sprintf("  %d. Description: %s█\n", i+1, m.input))
					b.WriteString(fmt.Sprintf("     Command: %s\n", cmd.Command))
				} else {
					b.WriteString(fmt.Sprintf("  %d. Description: %s\n", i+1, cmd.Description))
					b.WriteString(fmt.Sprintf("     Command: %s█\n", m.input))
				}
			} else {
				b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, cmd.Description))
				b.WriteString(fmt.Sprintf("     $ %s\n", cmd.Command))
			}
		}
		
		b.WriteString("\nPress Tab to switch field, Enter to add next command")
		b.WriteString("\nPress Ctrl+S to save")
	}

	b.WriteString("\n\nPress ESC to cancel")

	return listStyle.Render(b.String())
}

