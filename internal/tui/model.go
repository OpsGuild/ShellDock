package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shelldock/shelldock/internal/repo"
)

type view int

const (
	viewList view = iota
	viewDetail
	viewDeleteConfirm
	viewForm
)

type model struct {
	manager      *repo.Manager
	view         view
	sets         []string
	cursor       int
	scroll       int
	cmdSet       *repo.CommandSet
	detailScroll int
	form         *formModel
	deleteTarget string
	statusMsg    string
	width        int
	height       int
}

var (
	primaryColor = lipgloss.Color("#7C3AED")
	accentColor  = lipgloss.Color("#A78BFA")
	successColor = lipgloss.Color("#10B981")
	dangerColor  = lipgloss.Color("#EF4444")
	mutedColor   = lipgloss.Color("#6B7280")
	textColor    = lipgloss.Color("#F9FAFB")
	dimTextColor = lipgloss.Color("#9CA3AF")

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 2)

	statusStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	listItemStyle = lipgloss.NewStyle().
			Foreground(textColor).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				PaddingLeft(1)

	detailLabelStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true)

	detailValueStyle = lipgloss.NewStyle().
				Foreground(textColor)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#34D399"))

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	headerStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Underline(true)

	tagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#6D28D9")).
			Padding(0, 1)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimTextColor)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#818CF8")).
			Italic(true)
)

func newModel(manager *repo.Manager) *model {
	sets, _ := manager.GetLocalRepo().ListCommandSets()
	return &model{
		manager: manager,
		view:    viewList,
		sets:    sets,
		width:   80,
		height:  24,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) refreshSets() {
	sets, _ := m.manager.GetLocalRepo().ListCommandSets()
	m.sets = sets
	if m.cursor >= len(m.sets) {
		m.cursor = len(m.sets) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	switch m.view {
	case viewList:
		return m.updateList(msg)
	case viewDetail:
		return m.updateDetail(msg)
	case viewDeleteConfirm:
		return m.updateDelete(msg)
	case viewForm:
		return m.updateForm(msg)
	}
	return m, nil
}

func (m *model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.sets)-1 {
				m.cursor++
				maxVisible := m.maxListVisible()
				if m.cursor >= m.scroll+maxVisible {
					m.scroll = m.cursor - maxVisible + 1
				}
			}
		case "enter":
			if len(m.sets) > 0 {
				name := m.sets[m.cursor]
				cmdSet, err := m.manager.GetLocalRepo().GetCommandSet(name, "")
				if err == nil {
					m.cmdSet = cmdSet
					m.detailScroll = 0
					m.view = viewDetail
				} else {
					m.statusMsg = fmt.Sprintf("Error: %v", err)
				}
			}
		case "n":
			m.form = newFormModel(m.manager, nil, m.width, m.height)
			m.view = viewForm
			m.statusMsg = ""
		case "e":
			if len(m.sets) > 0 {
				name := m.sets[m.cursor]
				cmdSet, err := m.manager.GetLocalRepo().GetCommandSet(name, "")
				if err == nil {
					m.form = newFormModel(m.manager, cmdSet, m.width, m.height)
					m.view = viewForm
					m.statusMsg = ""
				}
			}
		case "d":
			if len(m.sets) > 0 {
				m.deleteTarget = m.sets[m.cursor]
				m.view = viewDeleteConfirm
			}
		}
	}
	return m, nil
}

func (m *model) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.view = viewList
			m.cmdSet = nil
			m.detailScroll = 0
		case "e":
			if m.cmdSet != nil {
				m.form = newFormModel(m.manager, m.cmdSet, m.width, m.height)
				m.view = viewForm
				m.statusMsg = ""
			}
		case "d":
			if m.cmdSet != nil {
				m.deleteTarget = m.cmdSet.Name
				m.view = viewDeleteConfirm
			}
		case "up", "k":
			if m.detailScroll > 0 {
				m.detailScroll--
			}
		case "down", "j":
			m.detailScroll++
		}
	}
	return m, nil
}

func (m *model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			err := m.manager.GetLocalRepo().DeleteCommandSet(m.deleteTarget)
			if err != nil {
				m.statusMsg = errorStyle.Render(fmt.Sprintf("Failed to delete: %v", err))
			} else {
				m.statusMsg = statusStyle.Render(fmt.Sprintf("Deleted '%s'", m.deleteTarget))
			}
			m.refreshSets()
			m.cmdSet = nil
			m.view = viewList
			m.deleteTarget = ""
		case "n", "N", "esc", "q":
			m.view = viewList
			if m.cmdSet != nil {
				m.view = viewDetail
			}
			m.deleteTarget = ""
		}
	}
	return m, nil
}

func (m *model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form == nil {
		m.view = viewList
		return m, nil
	}

	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)

	if m.form != nil && m.form.done {
		if m.form.saved {
			m.statusMsg = statusStyle.Render(fmt.Sprintf("Saved '%s'", m.form.cmdSet.Name))
		}
		m.form = nil
		m.refreshSets()
		m.view = viewList
	} else if m.form != nil && m.form.cancelled {
		m.form = nil
		m.view = viewList
	}

	return m, cmd
}

func (m *model) maxListVisible() int {
	return m.height - 10
}

func (m *model) View() string {
	switch m.view {
	case viewList:
		return m.viewList()
	case viewDetail:
		return m.viewDetail()
	case viewDeleteConfirm:
		return m.viewDeleteConfirm()
	case viewForm:
		if m.form != nil {
			return m.form.View()
		}
		return ""
	}
	return ""
}

func (m *model) viewList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" 📦 ShellDock "))
	b.WriteString("  ")
	b.WriteString(dimStyle.Render("Local Command Sets"))
	b.WriteString("\n\n")

	if len(m.sets) == 0 {
		b.WriteString(dimStyle.Render("  No command sets found in local repository."))
		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("  Command sets are collections of shell commands you can run together."))
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  They are stored as YAML files in ~/.shelldock/"))
		b.WriteString("\n\n")
		b.WriteString(dimStyle.Render("  Press 'n' to create your first command set."))
		b.WriteString("\n")
	} else {
		maxVisible := m.maxListVisible()
		if maxVisible < 3 {
			maxVisible = 3
		}

		end := m.scroll + maxVisible
		if end > len(m.sets) {
			end = len(m.sets)
		}

		if m.scroll > 0 {
			b.WriteString(dimStyle.Render("  ▲ more"))
			b.WriteString("\n")
		}

		for i := m.scroll; i < end; i++ {
			name := m.sets[i]
			if i == m.cursor {
				b.WriteString(selectedItemStyle.Render(fmt.Sprintf("❯ %s", name)))
			} else {
				b.WriteString(listItemStyle.Render(fmt.Sprintf("  %s", name)))
			}
			b.WriteString("\n")
		}

		if end < len(m.sets) {
			b.WriteString(dimStyle.Render("  ▼ more"))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(dimStyle.Render(fmt.Sprintf("  %d command set(s)", len(m.sets))))
		b.WriteString("\n")
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString("  " + m.statusMsg)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓ navigate • enter view details • n new • e edit • d delete • q quit"))
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  These are your local command sets in ~/.shelldock/"))

	return b.String()
}

func (m *model) viewDetail() string {
	if m.cmdSet == nil {
		return ""
	}

	var b strings.Builder
	cs := m.cmdSet

	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s ", cs.Name)))
	b.WriteString("\n\n")

	b.WriteString(detailLabelStyle.Render("Description  "))
	b.WriteString(detailValueStyle.Render(cs.Description))
	b.WriteString("\n")

	versionDisplay := cs.Version
	if cs.Tag != "" {
		versionDisplay = cs.Version + " @" + cs.Tag
	}
	b.WriteString(detailLabelStyle.Render("Version      "))
	b.WriteString(tagStyle.Render(versionDisplay))
	b.WriteString("\n\n")

	b.WriteString(headerStyle.Render("Steps"))
	b.WriteString("  ")
	b.WriteString(hintStyle.Render(fmt.Sprintf("%d step(s) — commands run in order from top to bottom", len(cs.Commands))))
	b.WriteString("\n\n")

	lines := []string{}

	for i, cmd := range cs.Commands {
		lines = append(lines, fmt.Sprintf("  %s %s",
			lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render(fmt.Sprintf("%d.", i+1)),
			detailValueStyle.Render(cmd.Description),
		))

		if cmd.Command != "" {
			lines = append(lines, fmt.Sprintf("     %s %s",
				dimStyle.Render("$"),
				commandStyle.Render(cmd.Command),
			))
		}

		if len(cmd.Platforms) > 0 {
			lines = append(lines, fmt.Sprintf("     %s", dimStyle.Render("platforms:")))
			for plat, platCmd := range cmd.Platforms {
				lines = append(lines, fmt.Sprintf("       %s %s",
					lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24")).Render(plat+":"),
					commandStyle.Render(platCmd),
				))
			}
		}

		if cmd.SkipOnError {
			lines = append(lines, fmt.Sprintf("     %s",
				dimStyle.Render("skip_on_error: true"),
			))
		}

		if len(cmd.Args) > 0 {
			lines = append(lines, fmt.Sprintf("     %s", dimStyle.Render("args:")))
			for _, arg := range cmd.Args {
				parts := []string{arg.Name}
				if arg.Required {
					parts = append(parts, "required")
				}
				if arg.Default != "" {
					parts = append(parts, fmt.Sprintf("default=%s", arg.Default))
				}
				if arg.Prompt != "" {
					parts = append(parts, fmt.Sprintf("prompt=%q", arg.Prompt))
				}
				lines = append(lines, fmt.Sprintf("       %s",
					dimStyle.Render(strings.Join(parts, " • ")),
				))
			}
		}

		lines = append(lines, "")
	}

	if len(cs.Commands) == 0 {
		lines = append(lines, dimStyle.Render("  No commands defined."))
	}

	maxLines := m.height - 10
	if maxLines < 3 {
		maxLines = 3
	}

	if m.detailScroll >= len(lines) {
		m.detailScroll = len(lines) - 1
	}
	if m.detailScroll < 0 {
		m.detailScroll = 0
	}

	endLine := m.detailScroll + maxLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	if m.detailScroll > 0 {
		b.WriteString(dimStyle.Render("  ▲ scroll up"))
		b.WriteString("\n")
	}

	for i := m.detailScroll; i < endLine; i++ {
		b.WriteString(lines[i])
		b.WriteString("\n")
	}

	if endLine < len(lines) {
		b.WriteString(dimStyle.Render("  ▼ scroll down"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓ scroll • e edit this command set • d delete • esc back to list"))
	b.WriteString("\n")
	b.WriteString(hintStyle.Render("  Press 'e' to modify steps, add platforms, arguments, and more."))

	return b.String()
}

func (m *model) viewDeleteConfirm() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" ⚠ Confirm Delete "))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  Delete command set %s?\n\n",
		errorStyle.Render(fmt.Sprintf("'%s'", m.deleteTarget)),
	))
	b.WriteString(dimStyle.Render("  This will permanently remove the YAML file from ~/.shelldock/"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  This action cannot be undone."))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("  y confirm delete • n/esc cancel and go back"))

	return b.String()
}
