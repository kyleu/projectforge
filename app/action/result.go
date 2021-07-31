package action

import (
	"fmt"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Result struct {
	ID       string           `json:"id"`
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
	msg := fmt.Sprintf("%+v", err)
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

func (r *Result) AddLog(msg string, args ...interface{}) {
	ret := fmt.Sprintf(msg, args...)
	if r.logger != nil {
		r.logger.Info(ret)
	}
	r.Logs = append(r.Logs, ret)
}

func (r *Result) AsError() error {
	if len(r.Errors) == 0 {
		return nil
	}
	var ret error
	for _, e := range r.Errors {
		if ret == nil {
			ret = errors.New(e)
		} else {
			ret = errors.Wrap(ret, e)
		}
	}
	return ret
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

type ResultContext struct {
	Prj *project.Project `json:"prj"`
	Cfg util.ValueMap `json:"cfg"`
	Res *Result `json:"res"`
}

func (c *ResultContext) Status() string {
	if c.Res == nil {
		return "Missing!"
	}
	fileCount := 0
	for _, m := range c.Res.Modules {
		for _, d := range m.Diffs {
			if d.Status != diff.StatusSkipped {
				fileCount++
			}
		}
	}
  if fileCount == 0 {
		return "ok"
	}
	return fmt.Sprintf("%d changes", fileCount)
}
