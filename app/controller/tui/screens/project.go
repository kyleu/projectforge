package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"projectforge.dev/projectforge/app/controller/tui/components"
	"projectforge.dev/projectforge/app/controller/tui/layout"
	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/controller/tui/style"
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

const (
	dataInputActive       = "input_active"
	dataInputChoiceKey    = "input_choice_key"
	dataInputChoiceTitle  = "input_choice_title"
	dataBuildCollapsed    = "build_collapsed"
	dataGitCollapsed      = "git_collapsed"
	dataListOffset        = "list_offset"
	dataViewportHeight    = "viewport_height"
	keyGroupBuild         = "group:build"
	keyGroupGit           = "group:git"
	prefixBuildSubActions = "action:build:"
	maxProjectRows        = 24
)

type ProjectScreen struct{}

type actionChoice struct {
	key         string
	title       string
	description string
	cfg         util.ValueMap
	runnable    bool
}

func NewProjectScreen() *ProjectScreen {
	return &ProjectScreen{}
}

func (s *ProjectScreen) Key() string {
	return KeyProject
}

func (s *ProjectScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	prj := selectedProject(ts, ps)
	if prj == nil {
		ps.Title = "Project"
		ps.SetStatus("Project not found")
		return nil
	}
	ps.Title = prj.Title()
	d := ps.EnsureData()
	d[dataBuildCollapsed] = true
	d[dataGitCollapsed] = true
	choices := s.visibleChoices(ps)
	if len(choices) == 0 {
		ps.Cursor = 0
	} else if ps.Cursor >= len(choices) {
		ps.Cursor = len(choices) - 1
	}
	ps.EnsureData()[dataListOffset] = 0
	ps.SetStatus("Choose an action for [%s]", prj.Key)
	return nil
}

func (s *ProjectScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	prj := selectedProject(ts, ps)
	allChoices := s.choices()
	choices := s.visibleChoices(ps)
	if len(choices) == 0 {
		ps.Cursor = 0
	} else if ps.Cursor >= len(choices) {
		ps.Cursor = len(choices) - 1
	}
	s.ensureViewportState(ps, len(choices))
	syncPromptChoice(ps, allChoices)

	if isPromptActive(ps) {
		return s.updatePrompt(ps, prj, msg)
	}

	switch m := msg.(type) {
	case tea.KeyMsg:
		switch m.String() {
		case "up", "k":
			s.moveCursor(ps, len(choices), -1)
		case "down", "j":
			s.moveCursor(ps, len(choices), 1)
		case "enter", "r", " ":
			if prj == nil || len(choices) == 0 {
				return mvc.Stay(), nil, nil
			}
			choice := choices[ps.Cursor]
			if s.isSection(choice.key) {
				s.toggleSection(ps, choice.key)
				post := s.visibleChoices(ps)
				if len(post) > 0 && ps.Cursor >= len(post) {
					ps.Cursor = len(post) - 1
				}
				s.ensureViewportState(ps, len(post))
				return mvc.Stay(), nil, nil
			}
			if !choice.runnable {
				ps.SetStatus("This row is not executable")
				return mvc.Stay(), nil, nil
			}
			if requiresInput(choice) {
				startInputPrompt(ps, choice)
				return mvc.Stay(), nil, nil
			}
			return mvc.Push(KeyResults, actionData(prj, choice, "", true)), nil, nil
		case "esc", "backspace", "b":
			return mvc.Pop(), nil, nil
		}
	}
	return mvc.Stay(), nil, nil
}

func (s *ProjectScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	styles := style.New(ts.Theme)
	prj := selectedProject(ts, ps)
	title := ps.Title
	if prj == nil {
		title = "Project Not Found"
	}

	choices := s.visibleChoices(ps)
	contentW, contentH, _ := mainPanelContentSize(styles.Panel, rects)
	contentH = min(contentH, maxProjectRows)
	ps.EnsureData()[dataViewportHeight] = contentH
	s.ensureViewportState(ps, len(choices))
	offset, end := s.viewportWindow(ps, len(choices))
	cursor := ps.Cursor - offset
	if cursor < 0 {
		cursor = 0
	}
	items := make(menu.Items, 0, end-offset)
	for idx := offset; idx < end; idx++ {
		c := choices[idx]
		items = append(items, &menu.Item{Key: c.key, Title: s.choiceTitle(ps, c), Description: c.description})
	}
	body := components.RenderMenuList(items, cursor, styles, contentW)
	if isPromptActive(ps) {
		body = lipgloss.JoinVertical(lipgloss.Left, body, styles.Muted.Render(renderPrompt(ps, contentW)))
	}

	panelOuterH := contentH + styles.Panel.GetVerticalFrameSize()
	panel := renderTextPanel(body, styles.Panel, rects.Main.W, panelOuterH)
	return renderScreenFrame(title, panel, styles, rects)
}

func (s *ProjectScreen) Help(_ *mvc.State, ps *mvc.PageState) HelpBindings {
	if isPromptActive(ps) {
		return HelpBindings{Short: []string{"type: message", "enter: continue", "tab: toggle dry-run", "esc: cancel"}}
	}
	return HelpBindings{Short: []string{"up/down: move", "enter: run/toggle", "b: back"}}
}

func (s *ProjectScreen) choices() []actionChoice {
	ret := make([]actionChoice, 0, 40)
	for _, t := range action.ProjectTypes {
		ret = append(ret, actionChoice{key: "action:" + t.Key, title: t.Title, description: t.Description, runnable: true})
	}
	for _, t := range []action.Type{action.TypeDebug, action.TypeRules, action.TypeSVG} {
		ret = append(ret, actionChoice{key: "action:" + t.Key, title: t.Title, description: t.Description, runnable: true})
	}
	ret = append(ret, actionChoice{key: keyGroupBuild, title: "Build", description: "Build phases", runnable: false})
	for _, b := range action.AllBuilds {
		if b == nil {
			continue
		}
		desc := b.Description
		if b.Expensive {
			desc = "[expensive] " + desc
		}
		ret = append(ret, actionChoice{
			key:         "action:" + action.TypeBuild.Key + ":" + b.Key,
			title:       "-> " + b.Title,
			description: desc,
			cfg:         util.ValueMap{"phase": b.Key},
			runnable:    true,
		})
	}
	ret = append(ret, actionChoice{key: keyGroupGit, title: "Git", description: "Repository actions", runnable: false})
	for _, ga := range []*git.Action{git.ActionStatus, git.ActionFetch, git.ActionPull, git.ActionPush, git.ActionCommit, git.ActionReset, git.ActionHistory, git.ActionMagic} {
		ret = append(ret, actionChoice{key: "git:" + ga.Key, title: "-> " + ga.Title, description: ga.Description, runnable: true})
	}
	return ret
}

func (s *ProjectScreen) visibleChoices(ps *mvc.PageState) []actionChoice {
	ret := make([]actionChoice, 0, 40)
	buildCollapsed := ps.EnsureData().GetBoolOpt(dataBuildCollapsed)
	gitCollapsed := ps.EnsureData().GetBoolOpt(dataGitCollapsed)
	for _, c := range s.choices() {
		switch {
		case c.key == keyGroupBuild:
			ret = append(ret, c)
		case strings.HasPrefix(c.key, prefixBuildSubActions):
			if !buildCollapsed {
				ret = append(ret, c)
			}
		case c.key == keyGroupGit:
			ret = append(ret, c)
		case strings.HasPrefix(c.key, "git:"):
			if !gitCollapsed {
				ret = append(ret, c)
			}
		default:
			ret = append(ret, c)
		}
	}
	return ret
}

func (s *ProjectScreen) choiceTitle(ps *mvc.PageState, c actionChoice) string {
	switch c.key {
	case keyGroupBuild:
		if ps.EnsureData().GetBoolOpt(dataBuildCollapsed) {
			return "[+] Build"
		}
		return "[-] Build"
	case keyGroupGit:
		if ps.EnsureData().GetBoolOpt(dataGitCollapsed) {
			return "[+] Git"
		}
		return "[-] Git"
	default:
		return c.title
	}
}

func (s *ProjectScreen) isSection(key string) bool {
	return key == keyGroupBuild || key == keyGroupGit
}

func (s *ProjectScreen) toggleSection(ps *mvc.PageState, key string) {
	d := ps.EnsureData()
	switch key {
	case keyGroupBuild:
		next := !d.GetBoolOpt(dataBuildCollapsed)
		d[dataBuildCollapsed] = next
		if next {
			ps.SetStatus("Collapsed [Build]")
		} else {
			ps.SetStatus("Expanded [Build]")
		}
	case keyGroupGit:
		next := !d.GetBoolOpt(dataGitCollapsed)
		d[dataGitCollapsed] = next
		if next {
			ps.SetStatus("Collapsed [Git]")
		} else {
			ps.SetStatus("Expanded [Git]")
		}
	}
}

func (s *ProjectScreen) updatePrompt(ps *mvc.PageState, prj *project.Project, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	selected := promptChoice(ps)
	d := ps.EnsureData()

	k, isKey := msg.(tea.KeyMsg)
	if !isKey {
		return mvc.Stay(), nil, nil
	}

	switch k.String() {
	case "esc":
		stopInputPrompt(ps)
		return mvc.Stay(), nil, nil
	case "backspace":
		v := d.GetStringOpt(dataInputMessage)
		r := []rune(v)
		if len(r) > 0 {
			d[dataInputMessage] = string(r[:len(r)-1])
		}
		return mvc.Stay(), nil, nil
	case "tab":
		if selected.key == "git:"+git.ActionMagic.Key {
			d[dataInputDryRun] = !d.GetBoolOpt(dataInputDryRun)
		}
		return mvc.Stay(), nil, nil
	case "enter":
		inputMessage := strings.TrimSpace(d.GetStringOpt(dataInputMessage))
		if inputMessage == "" {
			ps.SetStatus("Message is required")
			return mvc.Stay(), nil, nil
		}
		dryRun := d.GetBoolOpt(dataInputDryRun)
		stopInputPrompt(ps)
		return mvc.Push(KeyResults, actionData(prj, selected, inputMessage, dryRun)), nil, nil
	}

	if len(k.Runes) > 0 {
		d[dataInputMessage] = d.GetStringOpt(dataInputMessage) + string(k.Runes)
	}
	return mvc.Stay(), nil, nil
}

func requiresInput(c actionChoice) bool {
	return c.key == "git:"+git.ActionCommit.Key || c.key == "git:"+git.ActionMagic.Key
}

func isPromptActive(ps *mvc.PageState) bool {
	return ps.EnsureData().GetBoolOpt(dataInputActive)
}

func startInputPrompt(ps *mvc.PageState, c actionChoice) {
	d := ps.EnsureData()
	d[dataInputActive] = true
	d[dataInputChoiceKey] = c.key
	d[dataInputChoiceTitle] = c.title
	d[dataInputMessage] = defaultInputMessage(c)
	d[dataInputDryRun] = true
	ps.SetStatus("Enter values for [%s]", c.title)
}

func stopInputPrompt(ps *mvc.PageState) {
	d := ps.EnsureData()
	delete(d, dataInputActive)
	delete(d, dataInputChoiceKey)
	delete(d, dataInputChoiceTitle)
	delete(d, dataInputMessage)
	delete(d, dataInputDryRun)
	ps.SetStatus("Choose an action")
}

func promptChoice(ps *mvc.PageState) actionChoice {
	d := ps.EnsureData()
	return actionChoice{key: d.GetStringOpt(dataInputChoiceKey), title: d.GetStringOpt(dataInputChoiceTitle)}
}

func renderPrompt(ps *mvc.PageState, width int) string {
	d := ps.EnsureData()
	choice := promptChoice(ps)
	text := fmt.Sprintf("%s message: %s", choice.title, d.GetStringOpt(dataInputMessage))
	if choice.key == "git:"+git.ActionMagic.Key {
		mode := "on"
		if !d.GetBoolOpt(dataInputDryRun) {
			mode = "off"
		}
		text += fmt.Sprintf(" | dry-run: %s (tab to toggle)", mode)
	}
	return truncateLine(singleLine(text), width)
}

func syncPromptChoice(ps *mvc.PageState, choices []actionChoice) {
	if !isPromptActive(ps) {
		return
	}
	selected := promptChoice(ps)
	for _, c := range choices {
		if c.key == selected.key {
			ps.EnsureData()[dataInputChoiceTitle] = c.title
			return
		}
	}
	stopInputPrompt(ps)
}

func defaultInputMessage(c actionChoice) string {
	if c.key == "git:"+git.ActionMagic.Key {
		return "Project Forge TUI magic"
	}
	return "Project Forge TUI commit"
}

func (s *ProjectScreen) moveCursor(ps *mvc.PageState, count int, delta int) {
	if count == 0 || delta == 0 {
		return
	}
	next := ps.Cursor + delta
	if next < 0 {
		next = 0
	}
	if next >= count {
		next = count - 1
	}
	ps.Cursor = next
	s.ensureViewportState(ps, count)
}

func (s *ProjectScreen) ensureViewportState(ps *mvc.PageState, count int) {
	d := ps.EnsureData()
	if count <= 0 {
		d[dataListOffset] = 0
		return
	}
	height := d.GetIntOpt(dataViewportHeight)
	if height < 1 {
		height = 10
	}
	offset := d.GetIntOpt(dataListOffset)
	if offset < 0 {
		offset = 0
	}
	maxOffset := max(0, count-height)
	if offset > maxOffset {
		offset = maxOffset
	}
	if ps.Cursor < offset {
		offset = ps.Cursor
	}
	if ps.Cursor >= offset+height {
		offset = ps.Cursor - height + 1
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	d[dataListOffset] = offset
}

func (s *ProjectScreen) viewportWindow(ps *mvc.PageState, count int) (int, int) {
	if count <= 0 {
		return 0, 0
	}
	d := ps.EnsureData()
	offset := d.GetIntOpt(dataListOffset)
	height := d.GetIntOpt(dataViewportHeight)
	if height < 1 {
		height = 10
	}
	if offset < 0 {
		offset = 0
	}
	if offset > count-1 {
		offset = count - 1
	}
	end := min(count, offset+height)
	return offset, end
}

func renderLines(lines []string, width int) string {
	ret := make([]string, 0, len(lines))
	for _, line := range lines {
		ret = append(ret, truncateLine(singleLine(line), width))
	}
	return strings.Join(ret, "\n")
}

func truncateLine(s string, width int) string {
	if width < 1 {
		return ""
	}
	r := []rune(s)
	if len(r) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(r[:width-1]) + "…"
}

func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func runeLen(s string) int {
	return len([]rune(s))
}
