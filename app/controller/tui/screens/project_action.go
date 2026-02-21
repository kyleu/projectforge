package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"projectforge.dev/projectforge/app/controller/tui/mvc"
	"projectforge.dev/projectforge/app/lib/git"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
)

const (
	dataProjectKey     = "project"
	dataActionKey      = "action_key"
	dataActionTitle    = "action_title"
	dataActionDesc     = "action_desc"
	dataActionCfg      = "action_cfg"
	dataInputMessage   = "input_message"
	dataInputDryRun    = "input_dry_run"
	dataResultLines    = "result"
	dataResultTitle    = "result_title"
	dataResultComplete = "result_complete"
)

type projectResultMsg struct {
	title string
	lines []string
	err   error
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
		if prj == nil {
			return projectResultMsg{err: fmt.Errorf("project not found")}
		}
		if strings.HasPrefix(c.key, "action:") {
			actionKey := strings.TrimPrefix(c.key, "action:")
			if idx := strings.Index(actionKey, ":"); idx > -1 {
				actionKey = actionKey[:idx]
			}
			act := action.TypeFromString(actionKey)
			cfg := c.cfg
			if cfg == nil {
				cfg = util.ValueMap{}
			}
			params := &action.Params{
				ProjectKey: prj.Key,
				T:          act,
				Cfg:        cfg,
				MSvc:       ts.App.Services.Modules,
				PSvc:       ts.App.Services.Projects,
				ESvc:       ts.App.Services.Export,
				XSvc:       ts.App.Services.Exec,
				Logger:     logger,
			}
			res := action.Apply(ctx, params)
			lines := []string{fmt.Sprintf("status: %s", res.Status)}
			lines = append(lines, res.Logs...)
			for _, modRes := range res.Modules {
				if modRes == nil {
					continue
				}
				changes := len(modRes.DiffsFiltered(false))
				lines = append(lines, fmt.Sprintf("module status: %s (%s)", modRes.Status, util.StringPlural(changes, "change")))
				for _, d := range modRes.DiffsFiltered(false) {
					if d == nil || d.Status == nil {
						continue
					}
					lines = append(lines, fmt.Sprintf("  - %s: %s", d.Status.Key, d.Path))
				}
			}
			lines = append(lines, res.Errors...)
			if len(lines) == 1 {
				lines = append(lines, util.OK)
			}
			return projectResultMsg{title: fmt.Sprintf("Completed [%s]", act.Title), lines: lines}
		}

		gs := git.NewService(prj.Key, prj.Path)
		key := strings.TrimPrefix(c.key, "git:")
		var (
			res *git.Result
			err error
		)
		switch key {
		case git.ActionStatus.Key:
			res, err = gs.Status(ctx, logger)
		case git.ActionFetch.Key:
			res, err = gs.Fetch(ctx, logger)
		case git.ActionPull.Key:
			res, err = gs.Pull(ctx, logger)
		case git.ActionPush.Key:
			res, err = gs.Push(ctx, logger)
		case git.ActionCommit.Key:
			if message == "" {
				message = "Project Forge TUI commit"
			}
			res, err = gs.Commit(ctx, message, logger)
		case git.ActionReset.Key:
			res, err = gs.Reset(ctx, logger)
		case git.ActionHistory.Key:
			res, err = gs.History(ctx, &git.HistoryArgs{Limit: 25}, logger)
		case git.ActionMagic.Key:
			if message == "" {
				message = "Project Forge TUI magic"
			}
			res, err = gs.Magic(ctx, message, dryRun, logger)
		default:
			err = fmt.Errorf("unknown git action [%s]", key)
		}
		if err != nil {
			return projectResultMsg{err: err}
		}
		if res == nil {
			return projectResultMsg{err: fmt.Errorf("no result returned")}
		}
		lines := []string{fmt.Sprintf("status: %s", res.Status)}
		for _, k := range res.CleanData().Keys() {
			lines = append(lines, fmt.Sprintf("%s: %v", k, res.CleanData()[k]))
		}
		if res.Error != "" {
			lines = append(lines, "error: "+res.Error)
		}
		if len(lines) == 1 {
			lines = append(lines, util.OK)
		}
		return projectResultMsg{title: fmt.Sprintf("Completed [%s]", c.title), lines: lines}
	}
}

func selectedProject(ts *mvc.State, ps *mvc.PageState) *project.Project {
	if ts == nil || ts.App == nil || ts.App.Services == nil || ts.App.Services.Projects == nil {
		return nil
	}
	key := ps.EnsureData().GetStringOpt(dataProjectKey)
	if key == "" {
		return ts.App.Services.Projects.Default()
	}
	prj, err := ts.App.Services.Projects.Get(key)
	if err != nil {
		return nil
	}
	return prj
}
