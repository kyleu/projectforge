package action

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Result struct {
	Project  *project.Project `json:"project"`
	Action   Type             `json:"action"`
	Status   string           `json:"status"`
	Args     util.ValueMap    `json:"args,omitempty"`
	Data     any              `json:"data,omitempty"`
	Modules  module.Results   `json:"modules,omitempty"`
	Logs     []string         `json:"logs,omitempty"`
	Errors   []string         `json:"errors,omitempty"`
	Duration int              `json:"duration,omitempty"`
	logger   util.Logger
}

func newResult(act Type, prj *project.Project, cfg util.ValueMap, logger util.Logger) *Result {
	return &Result{Project: prj, Action: act, Args: cfg, Status: "OK", logger: logger}
}

func (r *Result) WithError(err error) *Result {
	msg := "error encountered"
	if err != nil {
		msg = err.Error()
	}
	if r.logger != nil {
		r.logger.Warnf("action error: %+v", err.Error())
	}
	r.Status = util.KeyError
	r.Errors = append(r.Errors, msg)
	return r
}

func (r *Result) AddDebug(msg string, args ...any) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Debug(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) AddLog(msg string, args ...any) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Info(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) AddWarn(msg string, args ...any) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Warn(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) Merge(tgt *Result) *Result {
	status := r.Status
	if status == "" {
		status = tgt.Status
	}
	logger := r.logger
	if logger == nil {
		logger = tgt.logger
	}
	return &Result{
		Status:   status,
		Args:     r.Args.Merge(tgt.Args),
		Modules:  append(append(module.Results{}, r.Modules...), tgt.Modules...),
		Logs:     append(append([]string{}, r.Logs...), tgt.Logs...),
		Errors:   append(append([]string{}, r.Errors...), tgt.Errors...),
		Duration: r.Duration + tgt.Duration,
		logger:   logger,
	}
}

func (r *Result) Title() string {
	return r.Action.Title
}

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *Result) AsError() error {
	if r.HasErrors() {
		return errors.New(strings.Join(r.Errors, "; "))
	}
	return nil
}

func (r *Result) StatusLog() string {
	fileCount := 0
	lo.ForEach(r.Modules, func(m *module.Result, _ int) {
		lo.ForEach(m.Diffs, func(d *diff.Diff, _ int) {
			if d.Status != diff.StatusSkipped {
				fileCount++
			}
		})
	})
	if fileCount == 0 {
		return "<em>ok</em>"
	}
	return util.StringPlural(fileCount, "change")
}

type ResultContext struct {
	Prj *project.Project `json:"prj,omitempty"`
	Cfg util.ValueMap    `json:"cfg,omitempty"`
	Res *Result          `json:"res,omitempty"`
}

func (c *ResultContext) Status() string {
	if c.Res == nil {
		return "<strong>missing</strong>"
	}
	return c.Res.StatusLog()
}

func (c *ResultContext) Title() string {
	if c.Res == nil {
		return "Unknown"
	}
	return c.Res.Action.Title
}

type ResultContexts []*ResultContext

func (x ResultContexts) Errors() []string {
	return lo.FlatMap(x, func(c *ResultContext, _ int) []string {
		if c.Res == nil {
			return nil
		}
		return c.Res.Errors
	})
}

func (x ResultContexts) Title() string {
	if len(x) == 0 || x[0].Res == nil {
		return fmt.Sprintf("Unknown (%d results)", len(x))
	}
	return x[0].Res.Action.Title
}

func errorResult(err error, t Type, cfg util.ValueMap, logger util.Logger) *Result {
	return newResult(t, nil, cfg, logger).WithError(err)
}
