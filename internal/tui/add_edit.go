package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shelldock/shelldock/internal/repo"
)

type focus int

const (
	focusMeta focus = iota
	focusSteps
	focusStep
	focusPlat
	focusArg
)

type metaField int

const (
	metaName metaField = iota
	metaDesc
	metaVersion
)

type stepField int

const (
	stepDesc stepField = iota
	stepCommand
	stepSkip
	stepPlatforms
	stepArgs
)

type platField int

const (
	platKey platField = iota
	platVal
)

type argField int

const (
	argName argField = iota
	argPrompt
	argDefault
	argRequired
)

var knownPlatforms = []string{
	"ubuntu", "debian", "centos", "rhel", "fedora",
	"arch", "opensuse", "alpine", "amazon", "oracle",
	"linux", "darwin", "windows",
}

type platEntry struct {
	key   string
	value string
}

type argEntry struct {
	name     string
	prompt   string
	defValue string
	required bool
}

type formModel struct {
	manager   *repo.Manager
	cmdSet    *repo.CommandSet
	isEdit    bool
	done      bool
	saved     bool
	cancelled bool
	err       string
	width     int
	height    int

	focus focus

	metaFld   metaField
	metaInput string

	stepCursor int

	openStep  int
	stepFld   stepField
	stepInput string

	editPlats []platEntry
	platIdx   int
	platFld   platField
	platInput string

	editArgs []argEntry
	argIdx   int
	argFld   argField
	argInput string

	scroll int
}

func newFormModel(manager *repo.Manager, cmdSet *repo.CommandSet, width, height int) *formModel {
	isEdit := cmdSet != nil
	if cmdSet == nil {
		cmdSet = &repo.CommandSet{
			Name:     "",
			Version:  "v1",
			Commands: []repo.Command{},
		}
	}

	cs := &repo.CommandSet{
		Name:        cmdSet.Name,
		Description: cmdSet.Description,
		Version:     cmdSet.Version,
		Commands:    make([]repo.Command, len(cmdSet.Commands)),
	}
	copy(cs.Commands, cmdSet.Commands)

	fm := &formModel{
		manager:    manager,
		cmdSet:     cs,
		isEdit:     isEdit,
		width:      width,
		height:     height,
		focus:      focusMeta,
		metaFld:    metaName,
		metaInput:  cs.Name,
		openStep:   -1,
		stepCursor: 0,
	}
	return fm
}

func (m *formModel) Init() tea.Cmd { return nil }

func (m *formModel) Update(msg tea.Msg) (*formModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.cancelled = true
			return m, nil
		}
		switch m.focus {
		case focusMeta:
			return m.updateMeta(msg)
		case focusSteps:
			return m.updateStepList(msg)
		case focusStep:
			return m.updateStepEdit(msg)
		case focusPlat:
			return m.updatePlatEdit(msg)
		case focusArg:
			return m.updateArgEdit(msg)
		}
	}
	return m, nil
}

func (m *formModel) updateMeta(msg tea.KeyMsg) (*formModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.cancelled = true
	case "tab", "enter":
		return m.advanceMeta()
	case "shift+tab":
		return m.retreatMeta()
	case "backspace":
		if len(m.metaInput) > 0 {
			m.metaInput = m.metaInput[:len(m.metaInput)-1]
		}
	default:
		switch msg.Type {
		case tea.KeyRunes:
			m.metaInput += string(msg.Runes)
		case tea.KeySpace:
			m.metaInput += " "
		}
	}
	return m, nil
}

func (m *formModel) advanceMeta() (*formModel, tea.Cmd) {
	switch m.metaFld {
	case metaName:
		val := strings.TrimSpace(m.metaInput)
		if val == "" {
			m.err = "Name is required"
			return m, nil
		}
		m.cmdSet.Name = val
		m.err = ""
		m.metaFld = metaDesc
		m.metaInput = m.cmdSet.Description
	case metaDesc:
		m.cmdSet.Description = strings.TrimSpace(m.metaInput)
		m.err = ""
		m.metaFld = metaVersion
		m.metaInput = m.cmdSet.Version
		if m.metaInput == "" {
			m.metaInput = "v1"
		}
	case metaVersion:
		val := strings.TrimSpace(m.metaInput)
		if val == "" {
			val = "v1"
		}
		m.cmdSet.Version = val
		m.err = ""
		m.focus = focusSteps
		m.stepCursor = 0
	}
	return m, nil
}

func (m *formModel) retreatMeta() (*formModel, tea.Cmd) {
	switch m.metaFld {
	case metaDesc:
		m.cmdSet.Description = strings.TrimSpace(m.metaInput)
		m.metaFld = metaName
		m.metaInput = m.cmdSet.Name
	case metaVersion:
		m.cmdSet.Version = strings.TrimSpace(m.metaInput)
		if m.cmdSet.Version == "" {
			m.cmdSet.Version = "v1"
		}
		m.metaFld = metaDesc
		m.metaInput = m.cmdSet.Description
	}
	m.err = ""
	return m, nil
}

func (m *formModel) updateStepList(msg tea.KeyMsg) (*formModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.focus = focusMeta
		m.metaFld = metaVersion
		m.metaInput = m.cmdSet.Version
		m.err = ""
	case "ctrl+s":
		return m.save()
	case "n":
		newCmd := repo.Command{}
		m.cmdSet.Commands = append(m.cmdSet.Commands, newCmd)
		idx := len(m.cmdSet.Commands) - 1
		m.stepCursor = idx
		m.openStepAt(idx)
	case "e", "enter":
		if len(m.cmdSet.Commands) > 0 && m.stepCursor >= 0 && m.stepCursor < len(m.cmdSet.Commands) {
			m.openStepAt(m.stepCursor)
		}
	case "d":
		if len(m.cmdSet.Commands) > 0 && m.stepCursor >= 0 && m.stepCursor < len(m.cmdSet.Commands) {
			idx := m.stepCursor
			m.cmdSet.Commands = append(m.cmdSet.Commands[:idx], m.cmdSet.Commands[idx+1:]...)
			if m.stepCursor >= len(m.cmdSet.Commands) {
				m.stepCursor = len(m.cmdSet.Commands) - 1
			}
			if m.stepCursor < 0 {
				m.stepCursor = 0
			}
		}
	case "up", "k":
		if m.stepCursor > 0 {
			m.stepCursor--
		}
	case "down", "j":
		if m.stepCursor < len(m.cmdSet.Commands)-1 {
			m.stepCursor++
		}
	}
	return m, nil
}

func (m *formModel) openStepAt(idx int) {
	m.openStep = idx
	m.focus = focusStep
	m.stepFld = stepDesc
	m.stepInput = m.cmdSet.Commands[idx].Description
	m.err = ""
}

func (m *formModel) updateStepEdit(msg tea.KeyMsg) (*formModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.commitStepField()
		m.closeStep()
	case "tab", "enter":
		return m.advanceStepField()
	case "shift+tab":
		return m.retreatStepField()
	case "backspace":
		if len(m.stepInput) > 0 {
			m.stepInput = m.stepInput[:len(m.stepInput)-1]
		}
	default:
		switch msg.Type {
		case tea.KeyRunes:
			m.stepInput += string(msg.Runes)
		case tea.KeySpace:
			m.stepInput += " "
		}
	}
	return m, nil
}

func (m *formModel) commitStepField() {
	if m.openStep < 0 || m.openStep >= len(m.cmdSet.Commands) {
		return
	}
	cmd := &m.cmdSet.Commands[m.openStep]
	switch m.stepFld {
	case stepDesc:
		cmd.Description = strings.TrimSpace(m.stepInput)
	case stepCommand:
		cmd.Command = strings.TrimSpace(m.stepInput)
	case stepSkip:
		val := strings.TrimSpace(strings.ToLower(m.stepInput))
		cmd.SkipOnError = val == "true" || val == "yes" || val == "y"
	}
}

func (m *formModel) advanceStepField() (*formModel, tea.Cmd) {
	m.commitStepField()
	cmd := &m.cmdSet.Commands[m.openStep]
	switch m.stepFld {
	case stepDesc:
		m.stepFld = stepCommand
		m.stepInput = cmd.Command
	case stepCommand:
		m.stepFld = stepSkip
		if cmd.SkipOnError {
			m.stepInput = "true"
		} else {
			m.stepInput = "false"
		}
	case stepSkip:
		m.stepFld = stepPlatforms
		m.focus = focusPlat
		m.loadPlatforms()
	case stepPlatforms:
		m.focus = focusPlat
		m.loadPlatforms()
	case stepArgs:
		m.focus = focusArg
		m.loadArgs()
	}
	m.err = ""
	return m, nil
}

func (m *formModel) retreatStepField() (*formModel, tea.Cmd) {
	m.commitStepField()
	cmd := &m.cmdSet.Commands[m.openStep]
	switch m.stepFld {
	case stepDesc:
		m.closeStep()
		return m, nil
	case stepCommand:
		m.stepFld = stepDesc
		m.stepInput = cmd.Description
	case stepSkip:
		m.stepFld = stepCommand
		m.stepInput = cmd.Command
	case stepPlatforms:
		m.stepFld = stepSkip
		if cmd.SkipOnError {
			m.stepInput = "true"
		} else {
			m.stepInput = "false"
		}
	case stepArgs:
		m.stepFld = stepPlatforms
		m.focus = focusPlat
		m.loadPlatforms()
	}
	m.err = ""
	return m, nil
}

func (m *formModel) closeStep() {
	m.openStep = -1
	m.focus = focusSteps
	m.err = ""
}

func (m *formModel) loadPlatforms() {
	cmd := m.cmdSet.Commands[m.openStep]
	m.editPlats = []platEntry{}
	for k, v := range cmd.Platforms {
		m.editPlats = append(m.editPlats, platEntry{key: k, value: v})
	}
	m.platIdx = 0
	m.platFld = platKey
	m.platInput = ""
	if len(m.editPlats) > 0 {
		m.platInput = m.editPlats[0].key
	}
}

func (m *formModel) savePlatforms() {
	if m.openStep < 0 || m.openStep >= len(m.cmdSet.Commands) {
		return
	}
	cmd := &m.cmdSet.Commands[m.openStep]
	plats := make(map[string]string)
	for _, p := range m.editPlats {
		if p.key != "" && p.value != "" {
			plats[p.key] = p.value
		}
	}
	if len(plats) > 0 {
		cmd.Platforms = plats
	} else {
		cmd.Platforms = nil
	}
}

func (m *formModel) commitPlatField() {
	if m.platIdx < 0 || m.platIdx >= len(m.editPlats) {
		return
	}
	switch m.platFld {
	case platKey:
		m.editPlats[m.platIdx].key = strings.TrimSpace(m.platInput)
	case platVal:
		m.editPlats[m.platIdx].value = strings.TrimSpace(m.platInput)
	}
}

func (m *formModel) updatePlatEdit(msg tea.KeyMsg) (*formModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.commitPlatField()
		m.savePlatforms()
		m.stepFld = stepArgs
		m.focus = focusArg
		m.loadArgs()
		return m, nil
	case "ctrl+n":
		m.commitPlatField()
		m.editPlats = append(m.editPlats, platEntry{})
		m.platIdx = len(m.editPlats) - 1
		m.platFld = platKey
		m.platInput = ""
	case "ctrl+d":
		if len(m.editPlats) > 0 && m.platIdx < len(m.editPlats) {
			m.editPlats = append(m.editPlats[:m.platIdx], m.editPlats[m.platIdx+1:]...)
			if m.platIdx >= len(m.editPlats) {
				m.platIdx = len(m.editPlats) - 1
			}
			if m.platIdx < 0 {
				m.platIdx = 0
			}
			m.platInput = ""
			if len(m.editPlats) > 0 {
				if m.platFld == platKey {
					m.platInput = m.editPlats[m.platIdx].key
				} else {
					m.platInput = m.editPlats[m.platIdx].value
				}
			}
		}
	case "tab", "enter":
		m.commitPlatField()
		if m.platFld == platKey {
			m.platFld = platVal
			if m.platIdx >= 0 && m.platIdx < len(m.editPlats) {
				m.platInput = m.editPlats[m.platIdx].value
			}
		} else {
			if m.platIdx < len(m.editPlats)-1 {
				m.platIdx++
				m.platFld = platKey
				m.platInput = m.editPlats[m.platIdx].key
			} else {
				m.savePlatforms()
				m.stepFld = stepArgs
				m.focus = focusArg
				m.loadArgs()
				return m, nil
			}
		}
	case "shift+tab":
		m.commitPlatField()
		if m.platFld == platVal {
			m.platFld = platKey
			if m.platIdx >= 0 && m.platIdx < len(m.editPlats) {
				m.platInput = m.editPlats[m.platIdx].key
			}
		} else {
			if m.platIdx > 0 {
				m.platIdx--
				m.platFld = platVal
				m.platInput = m.editPlats[m.platIdx].value
			} else {
				m.savePlatforms()
				m.focus = focusStep
				m.stepFld = stepSkip
				cmd := m.cmdSet.Commands[m.openStep]
				if cmd.SkipOnError {
					m.stepInput = "true"
				} else {
					m.stepInput = "false"
				}
				return m, nil
			}
		}
	case "backspace":
		if len(m.platInput) > 0 {
			m.platInput = m.platInput[:len(m.platInput)-1]
		}
	default:
		switch msg.Type {
		case tea.KeyRunes:
			m.platInput += string(msg.Runes)
		case tea.KeySpace:
			m.platInput += " "
		}
	}
	return m, nil
}

func (m *formModel) loadArgs() {
	cmd := m.cmdSet.Commands[m.openStep]
	m.editArgs = []argEntry{}
	for _, a := range cmd.Args {
		m.editArgs = append(m.editArgs, argEntry{
			name:     a.Name,
			prompt:   a.Prompt,
			defValue: a.Default,
			required: a.Required,
		})
	}
	m.argIdx = 0
	m.argFld = argName
	m.argInput = ""
	if len(m.editArgs) > 0 {
		m.argInput = m.editArgs[0].name
	}
}

func (m *formModel) saveArgs() {
	if m.openStep < 0 || m.openStep >= len(m.cmdSet.Commands) {
		return
	}
	cmd := &m.cmdSet.Commands[m.openStep]
	var args []repo.ArgumentDef
	for _, a := range m.editArgs {
		if a.name != "" {
			args = append(args, repo.ArgumentDef{
				Name:     a.name,
				Prompt:   a.prompt,
				Default:  a.defValue,
				Required: a.required,
			})
		}
	}
	cmd.Args = args
}

func (m *formModel) commitArgField() {
	if m.argIdx < 0 || m.argIdx >= len(m.editArgs) {
		return
	}
	a := &m.editArgs[m.argIdx]
	switch m.argFld {
	case argName:
		a.name = strings.TrimSpace(m.argInput)
	case argPrompt:
		a.prompt = strings.TrimSpace(m.argInput)
	case argDefault:
		a.defValue = strings.TrimSpace(m.argInput)
	case argRequired:
		val := strings.TrimSpace(strings.ToLower(m.argInput))
		a.required = val == "true" || val == "yes" || val == "y"
	}
}

func (m *formModel) curArgVal() string {
	if m.argIdx < 0 || m.argIdx >= len(m.editArgs) {
		return ""
	}
	a := m.editArgs[m.argIdx]
	switch m.argFld {
	case argName:
		return a.name
	case argPrompt:
		return a.prompt
	case argDefault:
		return a.defValue
	case argRequired:
		if a.required {
			return "true"
		}
		return "false"
	}
	return ""
}

func (m *formModel) updateArgEdit(msg tea.KeyMsg) (*formModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.commitArgField()
		m.saveArgs()
		m.closeStep()
		return m, nil
	case "ctrl+n":
		m.commitArgField()
		m.editArgs = append(m.editArgs, argEntry{})
		m.argIdx = len(m.editArgs) - 1
		m.argFld = argName
		m.argInput = ""
	case "ctrl+d":
		if len(m.editArgs) > 0 && m.argIdx < len(m.editArgs) {
			m.editArgs = append(m.editArgs[:m.argIdx], m.editArgs[m.argIdx+1:]...)
			if m.argIdx >= len(m.editArgs) {
				m.argIdx = len(m.editArgs) - 1
			}
			if m.argIdx < 0 {
				m.argIdx = 0
			}
			m.argInput = ""
			if len(m.editArgs) > 0 {
				m.argInput = m.curArgVal()
			}
		}
	case "tab", "enter":
		m.commitArgField()
		return m.advanceArgField()
	case "shift+tab":
		m.commitArgField()
		return m.retreatArgField()
	case "backspace":
		if len(m.argInput) > 0 {
			m.argInput = m.argInput[:len(m.argInput)-1]
		}
	default:
		switch msg.Type {
		case tea.KeyRunes:
			m.argInput += string(msg.Runes)
		case tea.KeySpace:
			m.argInput += " "
		}
	}
	return m, nil
}

func (m *formModel) advanceArgField() (*formModel, tea.Cmd) {
	switch m.argFld {
	case argName:
		m.argFld = argPrompt
	case argPrompt:
		m.argFld = argDefault
	case argDefault:
		m.argFld = argRequired
	case argRequired:
		if m.argIdx < len(m.editArgs)-1 {
			m.argIdx++
			m.argFld = argName
		} else {
			m.saveArgs()
			m.closeStep()
			return m, nil
		}
	}
	m.argInput = m.curArgVal()
	return m, nil
}

func (m *formModel) retreatArgField() (*formModel, tea.Cmd) {
	switch m.argFld {
	case argName:
		if m.argIdx > 0 {
			m.argIdx--
			m.argFld = argRequired
		} else {
			m.saveArgs()
			m.stepFld = stepPlatforms
			m.focus = focusPlat
			m.loadPlatforms()
			return m, nil
		}
	case argPrompt:
		m.argFld = argName
	case argDefault:
		m.argFld = argPrompt
	case argRequired:
		m.argFld = argDefault
	}
	m.argInput = m.curArgVal()
	return m, nil
}

func (m *formModel) save() (*formModel, tea.Cmd) {
	if m.cmdSet.Name == "" {
		m.err = "Name is required"
		return m, nil
	}
	var validCmds []repo.Command
	for _, cmd := range m.cmdSet.Commands {
		if cmd.Description != "" || cmd.Command != "" || len(cmd.Platforms) > 0 {
			validCmds = append(validCmds, cmd)
		}
	}
	m.cmdSet.Commands = validCmds

	if err := m.manager.GetLocalRepo().SaveCommandSet(m.cmdSet, ""); err != nil {
		m.err = fmt.Sprintf("Save failed: %v", err)
		return m, nil
	}
	m.done = true
	m.saved = true
	return m, nil
}

func (m *formModel) View() string {
	lines := m.buildLines()

	activeLine := m.findActiveLine(lines)
	visible := m.height - 1
	if visible < 5 {
		visible = 5
	}

	margin := 3
	if activeLine >= 0 {
		if activeLine < m.scroll+margin {
			m.scroll = activeLine - margin
		}
		if activeLine > m.scroll+visible-margin-1 {
			m.scroll = activeLine - visible + margin + 1
		}
	}
	if m.scroll < 0 {
		m.scroll = 0
	}
	maxScroll := len(lines) - visible
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}

	end := m.scroll + visible
	if end > len(lines) {
		end = len(lines)
	}

	var b strings.Builder
	for i := m.scroll; i < end; i++ {
		b.WriteString(lines[i].text)
		b.WriteString("\n")
	}
	return b.String()
}

type vline struct {
	text   string
	active bool
}

func (m *formModel) findActiveLine(lines []vline) int {
	for i, l := range lines {
		if l.active {
			return i
		}
	}
	return -1
}

func (m *formModel) buildLines() []vline {
	var lines []vline
	add := func(s string) { lines = append(lines, vline{text: s}) }
	addActive := func(s string) { lines = append(lines, vline{text: s, active: true}) }

	title := "New Command Set"
	if m.isEdit {
		title = "Edit Command Set"
	}
	add(titleStyle.Render(fmt.Sprintf(" ✏  %s ", title)))
	add("")

	if m.err != "" {
		add("  " + errorStyle.Render("⚠ "+m.err))
		add("")
	}

	type mf struct {
		label string
		value string
		fld   metaField
		hint  string
	}
	metas := []mf{
		{"Name", m.cmdSet.Name, metaName, "unique identifier, e.g. docker, my-setup"},
		{"Description", m.cmdSet.Description, metaDesc, "what this command set does"},
		{"Version", m.cmdSet.Version, metaVersion, "version tag, e.g. v1, v2"},
	}

	for _, f := range metas {
		active := m.focus == focusMeta && m.metaFld == f.fld
		label := dimStyle.Render(fmt.Sprintf("  %-14s", f.label))
		if active {
			label = detailLabelStyle.Render(fmt.Sprintf("❯ %-14s", f.label))
		}

		if active {
			addActive(label + renderInput(m.metaInput))
			add(hintStyle.Render(fmt.Sprintf("                  %s", f.hint)))
		} else {
			v := f.value
			if v == "" {
				v = "—"
			}
			add(label + detailValueStyle.Render(v))
		}
	}

	add("")

	if m.focus == focusMeta {
		add(helpStyle.Render("  tab/enter next field • shift+tab previous field • esc cancel"))
		add("")
		return lines
	}

	sepLen := 60
	if m.width > 10 {
		sepLen = m.width - 10
	}
	if sepLen > 80 {
		sepLen = 80
	}
	add(dimStyle.Render("  " + strings.Repeat("─", sepLen)))
	add("")

	stepsLabel := headerStyle.Render("  Steps")
	stepsCount := dimStyle.Render(fmt.Sprintf("  %d step(s)", len(m.cmdSet.Commands)))
	add(stepsLabel + stepsCount)
	add("")

	if len(m.cmdSet.Commands) == 0 && m.focus == focusSteps {
		add(dimStyle.Render("  No steps defined yet."))
		add(hintStyle.Render("  Press 'n' to add your first step."))
		addActive("")
	}

	for i, cmd := range m.cmdSet.Commands {
		isOpen := m.openStep == i
		isCursor := m.stepCursor == i && m.focus == focusSteps

		prefix := "  "
		style := dimStyle
		if isCursor {
			prefix = "❯ "
			style = detailValueStyle
		} else if isOpen {
			prefix = "▼ "
			style = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
		}

		desc := cmd.Description
		if desc == "" {
			desc = "(no description)"
		}
		headerLine := style.Render(fmt.Sprintf("%s%d. %s", prefix, i+1, desc))

		if isCursor {
			addActive(headerLine)
		} else {
			add(headerLine)
		}

		if !isOpen {
			extras := []string{}
			if cmd.Command != "" {
				c := cmd.Command
				if len(c) > 40 {
					c = c[:40] + "…"
				}
				extras = append(extras, fmt.Sprintf("$ %s", c))
			}
			if len(cmd.Platforms) > 0 {
				extras = append(extras, fmt.Sprintf("%d platform(s)", len(cmd.Platforms)))
			}
			if len(cmd.Args) > 0 {
				extras = append(extras, fmt.Sprintf("%d arg(s)", len(cmd.Args)))
			}
			if cmd.SkipOnError {
				extras = append(extras, "skip_on_error")
			}
			if len(extras) > 0 {
				add(dimStyle.Render(fmt.Sprintf("     %s", strings.Join(extras, " • "))))
			}
			add("")
			continue
		}

		add("")

		type sf struct {
			label string
			value string
			fld   stepField
			hint  string
		}
		skipVal := "false"
		if cmd.SkipOnError {
			skipVal = "true"
		}
		stepFields := []sf{
			{"Description", cmd.Description, stepDesc, "what this step does, e.g. \"Install Docker\""},
			{"Command", cmd.Command, stepCommand, "the shell command to run, e.g. apt-get update"},
			{"Skip on error", skipVal, stepSkip, "true/false — continues even if this step fails"},
		}

		for _, f := range stepFields {
			active := m.focus == focusStep && m.stepFld == f.fld
			label := dimStyle.Render(fmt.Sprintf("     %-16s", f.label))
			if active {
				label = detailLabelStyle.Render(fmt.Sprintf("   ❯ %-16s", f.label))
			}
			if active {
				addActive(label + renderInput(m.stepInput))
				add(hintStyle.Render(fmt.Sprintf("                        %s", f.hint)))
			} else {
				v := f.value
				if v == "" {
					v = "—"
				}
				add(label + detailValueStyle.Render(v))
			}
		}

		add("")

		platActive := m.focus == focusPlat || (m.focus == focusStep && m.stepFld == stepPlatforms)
		platLabel := dimStyle.Render("     Platforms")
		if platActive {
			platLabel = lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render("   ❯ Platforms")
		}
		platCount := len(cmd.Platforms)

		if m.focus == focusPlat {
			add(platLabel + dimStyle.Render(fmt.Sprintf("  (%d)", len(m.editPlats))))
			add(hintStyle.Render("       Supported: " + strings.Join(knownPlatforms, ", ")))
			add("")

			if len(m.editPlats) == 0 {
				addActive(dimStyle.Render("       No platforms defined. Press ctrl+n to add one."))
			}

			for pi, p := range m.editPlats {
				isCur := pi == m.platIdx

				kActive := isCur && m.platFld == platKey
				vActive := isCur && m.platFld == platVal

				kPrefix := "       "
				if isCur {
					kPrefix = "     ❯ "
				}

				kLabel := dimStyle.Render("platform: ")
				if kActive {
					addActive(detailLabelStyle.Render(kPrefix) + kLabel + renderInput(m.platInput))
					add(hintStyle.Render(fmt.Sprintf("                type a platform name, e.g. %s", suggestPlatform(m.platInput))))
				} else {
					kStr := p.key
					if kStr == "" {
						kStr = "—"
					}
					add(dimStyle.Render(kPrefix) + kLabel + detailValueStyle.Render(kStr))
				}

				vLabel := dimStyle.Render("  command: ")
				if vActive {
					addActive("       " + vLabel + renderInput(m.platInput))
					add(hintStyle.Render("                the shell command to run on this platform"))
				} else {
					vStr := p.value
					if vStr == "" {
						vStr = "—"
					}
					add("       " + vLabel + detailValueStyle.Render(vStr))
				}
			}

			add("")
			add(helpStyle.Render("       tab/enter next • shift+tab back • ctrl+n add • ctrl+d remove"))
			add(helpStyle.Render("       esc done with platforms → move to arguments"))
		} else {
			if platCount > 0 {
				pKeys := []string{}
				for k := range cmd.Platforms {
					pKeys = append(pKeys, k)
				}
				add(platLabel + dimStyle.Render(fmt.Sprintf("  %d: %s", platCount, strings.Join(pKeys, ", "))))
			} else {
				add(platLabel + dimStyle.Render("  none"))
			}
		}

		add("")

		argActive := m.focus == focusArg || (m.focus == focusStep && m.stepFld == stepArgs)
		argLabel := dimStyle.Render("     Arguments")
		if argActive {
			argLabel = lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render("   ❯ Arguments")
		}

		if m.focus == focusArg {
			add(argLabel + dimStyle.Render(fmt.Sprintf("  (%d)", len(m.editArgs))))
			add(hintStyle.Render("       Use {{name}} placeholders in commands."))
			add("")

			if len(m.editArgs) == 0 {
				addActive(dimStyle.Render("       No arguments defined. Press ctrl+n to add one."))
			}

			for ai, a := range m.editArgs {
				isCur := ai == m.argIdx

				type af struct {
					label string
					value string
					fld   argField
					hint  string
				}
				reqStr := "false"
				if a.required {
					reqStr = "true"
				}
				aFields := []af{
					{"name", a.name, argName, "variable name for {{name}} placeholders"},
					{"prompt", a.prompt, argPrompt, "question shown to user"},
					{"default", a.defValue, argDefault, "value used if user presses enter without typing"},
					{"required", reqStr, argRequired, "true/false — if true, user must provide a value"},
				}

				for fi, f := range aFields {
					isAct := isCur && m.argFld == f.fld

					var prefix string
					if fi == 0 && isCur {
						prefix = "     ❯ "
					} else if fi == 0 {
						prefix = "       "
					} else {
						prefix = "         "
					}

					lbl := dimStyle.Render(fmt.Sprintf("%-12s", f.label+":"))
					if isAct {
						addActive(detailLabelStyle.Render(prefix) + lbl + renderInput(m.argInput))
						add(hintStyle.Render(fmt.Sprintf("                      %s", f.hint)))
					} else {
						v := f.value
						if v == "" {
							v = "—"
						}
						add(dimStyle.Render(prefix) + lbl + detailValueStyle.Render(v))
					}
				}
				if ai < len(m.editArgs)-1 {
					add("")
				}
			}

			add("")
			add(helpStyle.Render("       tab/enter next • shift+tab back • ctrl+n add • ctrl+d remove"))
			add(helpStyle.Render("       esc done with arguments → close step"))
		} else {
			argCount := len(cmd.Args)
			if argCount > 0 {
				aNames := []string{}
				for _, a := range cmd.Args {
					aNames = append(aNames, "{{"+a.Name+"}}")
				}
				add(argLabel + dimStyle.Render(fmt.Sprintf("  %d: %s", argCount, strings.Join(aNames, ", "))))
			} else {
				add(argLabel + dimStyle.Render("  none"))
			}
		}

		add("")

		if m.focus == focusStep {
			add(helpStyle.Render("     tab/enter next field • shift+tab previous • esc close step"))
		}

		add("")
	}

	add(dimStyle.Render("  " + strings.Repeat("─", sepLen)))
	add("")

	switch m.focus {
	case focusSteps:
		add(helpStyle.Render("  n add step • e/enter edit step • d remove step • ↑/↓ navigate"))
		add(helpStyle.Render("  ctrl+s save & exit • esc go back to metadata"))
	case focusStep, focusPlat, focusArg:
		add(helpStyle.Render("  Editing step — use tab/shift+tab to move between fields"))
		add(helpStyle.Render("  esc closes current section • ctrl+s saves from step list"))
	}

	return lines
}

func suggestPlatform(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return "ubuntu"
	}
	for _, p := range knownPlatforms {
		if strings.HasPrefix(p, input) {
			return p
		}
	}
	return input
}

func renderInput(value string) string {
	cursor := lipgloss.NewStyle().
		Background(lipgloss.Color("#7C3AED")).
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(" ")
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#F9FAFB")).Render(value) + cursor
}
