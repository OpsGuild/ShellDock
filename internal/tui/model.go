package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shelldock/shelldock/internal/repo"
)

type state int

const (
	stateList state = iota
	stateView
	stateEdit
	stateAdd
	stateDelete
)

type model struct {
	manager     *repo.Manager
	state       state
	commandSets []string
	selected    int
	cmdSet      *repo.CommandSet
	editingIdx  int
	input       string
	err         error
	addEdit     *addEditModel
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	listStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Margin(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))
)

func NewModel(manager *repo.Manager) *model {
	sets, _ := manager.GetLocalRepo().ListCommandSets()
	return &model{
		manager:     manager,
		state:       stateList,
		commandSets: sets,
		selected:    0,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.addEdit != nil && (m.state == stateAdd || m.state == stateEdit) {
		var cmd tea.Cmd
		m.addEdit, cmd = m.addEdit.Update(msg)
		if m.addEdit == nil {
			// User quit or saved
			m.addEdit = nil
			if m.state == stateAdd || m.state == stateEdit {
				// Refresh list
				sets, _ := m.manager.GetLocalRepo().ListCommandSets()
				m.commandSets = sets
				m.state = stateList
			}
			return m, cmd
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateView:
			return m.updateView(msg)
		case stateEdit:
			return m.updateEdit(msg)
		case stateAdd:
			return m.updateAdd(msg)
		}
	}
	return m, nil
}

func (m *model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case "down", "j":
		if m.selected < len(m.commandSets)-1 {
			m.selected++
		}
	case "enter":
		if len(m.commandSets) > 0 {
			name := m.commandSets[m.selected]
			cmdSet, err := m.manager.GetLocalRepo().GetCommandSet(name, "")
			if err == nil {
				m.cmdSet = cmdSet
				m.state = stateView
			}
		}
	case "a":
		m.state = stateAdd
		m.addEdit = newAddEditModel(m.manager, false, nil)
	case "d":
		if len(m.commandSets) > 0 {
			m.state = stateDelete
		}
	}
	return m, nil
}

func (m *model) updateView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = stateList
		m.cmdSet = nil
	case "e":
		if m.cmdSet != nil {
			m.state = stateEdit
			m.addEdit = newAddEditModel(m.manager, true, m.cmdSet)
		}
	}
	return m, nil
}

func (m *model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handled by addEdit model
	return m, nil
}

func (m *model) updateAdd(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handled by addEdit model
	return m, nil
}

func (m *model) View() string {
	switch m.state {
	case stateList:
		return m.viewList()
	case stateView:
		return m.viewCommandSet()
	case stateEdit:
		return m.viewEdit()
	case stateAdd:
		return m.viewAdd()
	case stateDelete:
		return m.viewDelete()
	default:
		return ""
	}
}

func (m *model) viewList() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(" ShellDock - Command Sets "))
	b.WriteString("\n\n")

	if len(m.commandSets) == 0 {
		b.WriteString("No command sets found. Press 'a' to add one.\n")
	} else {
		for i, name := range m.commandSets {
			style := normalStyle
			if i == m.selected {
				style = selectedStyle
				b.WriteString("▶ ")
			} else {
				b.WriteString("  ")
			}
			b.WriteString(style.Render(name))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString("Controls: ↑/↓ Navigate | Enter View | a Add | d Delete | q Quit\n")

	return listStyle.Render(b.String())
}

func (m *model) viewCommandSet() string {
	if m.cmdSet == nil {
		return "No command set selected"
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s ", m.cmdSet.Name)))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Description: %s\n", m.cmdSet.Description))
	b.WriteString(fmt.Sprintf("Version: %s\n\n", m.cmdSet.Version))
	b.WriteString("Commands:\n\n")

	for i, cmd := range m.cmdSet.Commands {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, cmd.Description))
		b.WriteString(fmt.Sprintf("   $ %s\n\n", cmd.Command))
	}

	b.WriteString("\nControls: e Edit | q/ESC Back\n")

	return listStyle.Render(b.String())
}

func (m *model) viewEdit() string {
	if m.addEdit != nil {
		return m.addEdit.View()
	}
	return listStyle.Render("Edit mode\n\nPress ESC to cancel")
}

func (m *model) viewAdd() string {
	if m.addEdit != nil {
		return m.addEdit.View()
	}
	return listStyle.Render("Add mode\n\nPress ESC to cancel")
}

func (m *model) viewDelete() string {
	if len(m.commandSets) == 0 {
		m.state = stateList
		return m.viewList()
	}

	name := m.commandSets[m.selected]
	m.manager.GetLocalRepo().DeleteCommandSet(name)
	sets, _ := m.manager.GetLocalRepo().ListCommandSets()
	m.commandSets = sets
	m.state = stateList
	return m.viewList()
}

