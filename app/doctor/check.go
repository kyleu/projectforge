package doctor

import (
	"context"
	"runtime"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

type (
	CheckFn func(ctx context.Context, r *Result, logger util.Logger) *Result
	SolveFn func(ctx context.Context, r *Result, logger util.Logger) *Result
)

type Check struct {
	Key       string   `json:"key"`
	Section   string   `json:"section"`
	Title     string   `json:"title"`
	Summary   string   `json:"summary,omitzero"`
	URL       string   `json:"url,omitzero"`
	UsedBy    string   `json:"usedBy,omitzero"`
	Modules   []string `json:"modules,omitempty"`
	Platforms []string `json:"platforms,omitempty"`
	Core      bool     `json:"core,omitzero"`
	Fn        CheckFn  `json:"-"`
	Solve     SolveFn  `json:"-"`
}

func (c *Check) Check(ctx context.Context, logger util.Logger) *Result {
	_, span, logger := telemetry.StartSpan(ctx, "doctor:check:"+c.Key, logger)
	defer span.Complete()

	if !c.Applies() {
		return nil
	}

	r := NewResult(c, c.Key, c.Title, c.Summary)
	timer := util.TimerStart()
	r = c.Fn(ctx, r, logger)
	r.Duration = timer.End()
	r = c.Solve(ctx, r, logger)
	return r
}

func (c *Check) Applies() bool {
	if len(c.Platforms) == 0 {
		return true
	}
	return lo.Contains(c.Platforms, runtime.GOOS) && (!lo.Contains(c.Platforms, "!"+runtime.GOOS))
}

type Checks []*Check

func (c Checks) Get(key string) *Check {
	return lo.FindOrElse(c, nil, func(x *Check) bool {
		return x.Key == key
	})
}

func (c Checks) Keys() []string {
	return lo.Map(c, func(c *Check, _ int) string {
		return c.Key
	})
}
