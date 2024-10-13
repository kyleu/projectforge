package task

import (
	"context"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/exec"
	"{{{ .Package }}}/app/util"
)

type TaskFn func(ctx context.Context, res *Result, logger util.Logger) error

type Task struct {
	Key           string          `json:"key"`
	Title         string          `json:"title,omitempty"`
	Category      string          `json:"category,omitempty"`
	Icon          string          `json:"icon,omitempty"`
	Description   string          `json:"description,omitempty"`
	Tags          []string        `json:"tags,omitempty"`
	Fields        util.FieldDescs `json:"fields,omitempty"`
	Dangerous     string          `json:"dangerous,omitempty"`
	WebURL        string          `json:"webURL,omitempty"`
	MaxConcurrent int             `json:"maxConcurrent,omitempty"`
	fns           []TaskFn
}

func NewTask(key string, title string, cat string, icon string, desc string, fns ...TaskFn) *Task {
	if title == "" {
		title = util.StringToTitle(key)
	}
	return &Task{Key: key, Title: title, Category: cat, Icon: icon, Description: desc, fns: fns}
}

func (t *Task) TitleSafe() string {
	if t == nil {
		return "<nil task>"
	}
	if t.Title == "" {
		return t.Key
	}
	return t.Title
}

func (t *Task) IconSafe() string {
	if t == nil || t.Icon == "" {
		return "star"
	}
	return t.Icon
}

func (t *Task) WebPath() string {
	if t.WebURL != "" {
		return t.WebURL
	}
	return "/admin/task/" + t.Key
}

func (t *Task) Clone() *Task {
	return &Task{
		Key: t.Key, Title: t.Title, Category: t.Category, Icon: t.Icon, Description: t.Description, Tags: t.Tags,
		Fields: t.Fields, Dangerous: t.Dangerous, MaxConcurrent: t.MaxConcurrent,
	}
}

func (t *Task) WithTags(tags []string) *Task {
	ret := t.Clone()
	ret.Tags = tags
	return ret
}

func (t *Task) Run(ctx context.Context, runner string, args util.ValueMap, logger util.Logger, fns ...exec.OutFn) *Result {
	ret := NewResult(t, runner, args, t.resultLogFn(logger, fns...))
	return t.RunWithResult(ctx, ret, logger)
}

func (t *Task) RunWithResult(ctx context.Context, res *Result, logger util.Logger) *Result {
	if len(t.fns) == 0 {
		res.Log("no work to do for task [%s]", t.TitleSafe())
		return res.Complete("OK")
	}
	for _, fn := range t.fns {
		err := fn(ctx, res, logger)
		if err != nil {
			return res.CompleteError(errors.Wrapf(err, "unable to run task [%s]", t.Key))
		}
	}
	return res.Complete("OK")
}

func (t *Task) resultLogFn(logger util.Logger, fns ...exec.OutFn) ResultLogFn {
	return func(key string, data any) {
		for _, fn := range fns {
			var b []byte
			if s, ok := data.(string); ok {
				b = []byte(s)
			} else {
				util.ToJSONBytes(data, true)
			}
			if err := fn(key, b); err != nil {
				logger.Warnf("error calling result function: %s", err.Error())
			}
		}
	}
}
