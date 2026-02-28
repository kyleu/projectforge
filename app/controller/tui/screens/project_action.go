package screens

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

const (
	dataProjectKey      = "project"
	dataActionKey       = "action_key"
	dataActionTitle     = "action_title"
	dataActionDesc      = "action_desc"
	dataActionCfg       = "action_cfg"
	dataInputMessage    = "input_message"
	dataInputDryRun     = "input_dry_run"
	dataResultLines     = "result"
	dataResultTitle     = "result_title"
	dataResultComplete  = "result_complete"
	dataResultChanges   = "result_changes"
	dataResultStatus    = "result_status"
	dataResultDiffPath  = "result_diff_path"
	dataResultDiffTag   = "result_diff_tag"
	dataResultDiffPatch = "result_diff_patch"
)

type projectResultMsg struct {
	title   string
	status  string
	lines   []string
	changes []resultChange
	err     error
}

type resultChange struct {
	StatusKey string
	Path      string
	Patch     string
}

func actionData(prj *project.Project, c actionChoice, message string, dryRun bool) util.ValueMap {
	ret := util.ValueMap{
		dataActionKey:   c.key,
		dataActionTitle: c.title,
		dataActionDesc:  c.description,
		dataInputDryRun: dryRun,
	}
	if c.cfg != nil {
		ret[dataActionCfg] = c.cfg.Clone()
	}
	if prj != nil {
		ret[dataProjectKey] = prj.Key
	}
	if message != "" {
		ret[dataInputMessage] = message
	}
	return ret
}

func actionChoiceFromData(data util.ValueMap) actionChoice {
	if data == nil {
		return actionChoice{}
	}
	cfg, _ := data.GetMap(dataActionCfg, true)
	return actionChoice{
		key:         data.GetStringOpt(dataActionKey),
		title:       data.GetStringOpt(dataActionTitle),
		description: data.GetStringOpt(dataActionDesc),
		cfg:         cfg,
		runnable:    true,
	}
}

func runProjectActionCmd(ts *mvc.State, ps *mvc.PageState, prj *project.Project, c actionChoice, message string, dryRun bool) tea.Cmd {
	ctx := ps.Context
	logger := ps.Logger
	if logger == nil {
		logger = ts.Logger
	}
	return func() tea.Msg {
		if strings.HasPrefix(c.key, "action:") {
			return runModuleAction(ctx, ts, prj, c, logger)
		}
		if prj == nil {
			return projectResultMsg{err: errors.New("project is nil")}
		}
		return runGitAction(ctx, prj, c, message, dryRun, logger)
	}
}

func runModuleAction(ctx context.Context, ts *mvc.State, prj *project.Project, c actionChoice, logger util.Logger) projectResultMsg {
	actionKey := strings.TrimPrefix(c.key, "action:")
	if idx := strings.Index(actionKey, ":"); idx > -1 {
		actionKey = actionKey[:idx]
	}
	act := action.TypeFromString(actionKey)
	cfg := c.cfg
	if cfg == nil {
		cfg = util.ValueMap{}
	}
	projectKey := ""
	if prj != nil {
		projectKey = prj.Key
	}
	params := &action.Params{
		ProjectKey: projectKey,
		T:          act,
		Cfg:        cfg,
		MSvc:       ts.App.Services.Modules,
		PSvc:       ts.App.Services.Projects,
		ESvc:       ts.App.Services.Export,
		XSvc:       ts.App.Services.Exec,
		CLI:        cfg.GetBoolOpt("cli"),
		Logger:     logger,
	}
	res := action.Apply(ctx, params)
	changes := collectResultChanges(res)
	lines := moduleResultLines(res, changes)
	return projectResultMsg{title: fmt.Sprintf("Completed [%s]", act.Title), status: res.Status, lines: lines, changes: changes}
}

func collectResultChanges(res *action.Result) []resultChange {
	changes := make([]resultChange, 0, 32)
	for _, modRes := range res.Modules {
		if modRes == nil {
			continue
		}
		for _, d := range modRes.DiffsFiltered(false) {
			if d == nil || d.Status == nil {
				continue
			}
			changes = append(changes, resultChange{StatusKey: d.Status.Key, Path: d.Path, Patch: d.Patch})
		}
	}
	return changes
}

func moduleResultLines(res *action.Result, changes []resultChange) []string {
	lines := []string{fmt.Sprintf("Status: %s", res.Status)}
	if len(res.Modules) > 0 {
		lines = append(lines, fmt.Sprintf("Changes: %s", util.StringPlural(len(changes), "change")))
		for _, ch := range changes {
			lines = append(lines, fmt.Sprintf("[%s] %s", resultChangeTag(ch.StatusKey), ch.Path))
		}
	} else {
		lines = append(lines, res.Logs...)
	}
	lines = append(lines, res.Errors...)
	if len(lines) == 1 {
		lines = append(lines, util.KeyOK)
	}
	return lines
}

func runGitAction(
	ctx context.Context,
	prj *project.Project,
	c actionChoice,
	message string,
	dryRun bool,
	logger util.Logger,
) projectResultMsg {
	gs := git.NewService(prj.Key, prj.Path)
	key := strings.TrimPrefix(c.key, "git:")
	res, err := executeGitAction(ctx, gs, key, message, dryRun, logger)
	if err != nil {
		return projectResultMsg{err: err}
	}
	if res == nil {
		return projectResultMsg{err: errors.Errorf("no git result from action [%s]", c.key)}
	}
	lines := []string{fmt.Sprintf("status: %s", res.Status)}
	for _, k := range res.CleanData().Keys() {
		lines = append(lines, fmt.Sprintf("%s: %v", k, res.CleanData()[k]))
	}
	if res.Error != "" {
		lines = append(lines, "error: "+res.Error)
	}
	if len(lines) == 1 {
		lines = append(lines, util.KeyOK)
	}
	return projectResultMsg{title: fmt.Sprintf("Completed [%s]", c.title), status: res.Status, lines: lines}
}

func executeGitAction(
	ctx context.Context,
	gs *git.Service,
	key string,
	message string,
	dryRun bool,
	logger util.Logger,
) (*git.Result, error) {
	switch key {
	case git.ActionStatus.Key:
		return gs.Status(ctx, logger)
	case git.ActionFetch.Key:
		return gs.Fetch(ctx, logger)
	case git.ActionPull.Key:
		return gs.Pull(ctx, logger)
	case git.ActionPush.Key:
		return gs.Push(ctx, logger)
	case git.ActionCommit.Key:
		return gs.Commit(ctx, message, logger)
	case git.ActionReset.Key:
		return gs.Reset(ctx, logger)
	case git.ActionHistory.Key:
		return gs.History(ctx, &git.HistoryArgs{Limit: 25}, logger)
	case git.ActionMagic.Key:
		if message == "" {
			message = "Project Forge TUI magic"
		}
		return gs.Magic(ctx, message, dryRun, logger)
	default:
		return nil, errors.Errorf("unknown git action [%s]", key)
	}
}

func resultChangeTag(statusKey string) string {
	switch statusKey {
	case diff.StatusDifferent.Key:
		return "M"
	case diff.StatusNew.Key:
		return "A"
	case diff.StatusMissing.Key:
		return "D"
	default:
		return "?"
	}
}

func selectedProject(ts *mvc.State, ps *mvc.PageState) *project.Project {
	if ts == nil || ts.App == nil || ts.App.Services == nil || ts.App.Services.Projects == nil {
		return nil
	}
	data := ps.EnsureData()
	rawKey, hasProjectKey := data[dataProjectKey]
	if !hasProjectKey {
		return ts.App.Services.Projects.Default()
	}
	key, _ := rawKey.(string)
	if key == "" {
		return nil
	}
	prj, err := ts.App.Services.Projects.Get(key)
	if err != nil {
		return nil
	}
	return prj
}
