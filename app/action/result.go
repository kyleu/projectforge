package action

import (
	"fmt"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

type Result struct {
	Project  *project.Project `json:"project"`
	Status   string           `json:"status"`
	Args     util.ValueMap    `json:"args,omitempty"`
	Modules  module.Results   `json:"modules,omitempty"`
	Logs     []string         `json:"logs,omitempty"`
	Errors   []string         `json:"errors,omitempty"`
	Duration int              `json:"duration,omitempty"`
	logger   *zap.SugaredLogger
}

func newResult(cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	return &Result{Args: cfg, Status: "OK", logger: logger}
}

func (r *Result) WithError(err error) *Result {
	msg := err.Error()
	if r.logger != nil {
		r.logger.Errorf("%+v", err)
	}
	r.Status = "error"
	r.Errors = append(r.Errors, msg)
	return r
}

func errorResult(err error, cfg util.ValueMap, logger *zap.SugaredLogger) *Result {
	return newResult(cfg, logger).WithError(err)
}

func (r *Result) AddDebug(msg string, args ...interface{}) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Debug(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) AddLog(msg string, args ...interface{}) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Info(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) AddWarn(msg string, args ...interface{}) {
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

func (r *Result) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *Result) StatusLog() string {
	fileCount := 0
	for _, m := range r.Modules {
		for _, d := range m.Diffs {
			if d.Status != diff.StatusSkipped {
				fileCount++
			}
		}
	}
	if fileCount == 0 {
		return "no changes"
	}
	return fmt.Sprintf("%d %s", fileCount, util.StringPluralMaybe("change", fileCount))
}

type ResultContext struct {
	Prj *project.Project `json:"prj,omitempty"`
	Cfg util.ValueMap    `json:"cfg,omitempty"`
	Res *Result          `json:"res,omitempty"`
}

func (c *ResultContext) Status() string {
	if c.Res == nil {
		return "Missing!"
	}
	return c.Res.StatusLog()
}

type ResultContexts []*ResultContext

func (x ResultContexts) Errors() []string {
	var ret []string
	for _, c := range x {
		if c.Res != nil {
			ret = append(ret, c.Res.Errors...)
		}
	}
	return ret
}
