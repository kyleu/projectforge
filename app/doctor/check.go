package doctor

import (
	"context"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

type (
	CheckFn func(ctx context.Context, r *Result, logger util.Logger) *Result
	SolveFn func(ctx context.Context, r *Result, logger util.Logger) *Result
)

type Check struct {
	Key     string   `json:"key"`
	Section string   `json:"section"`
	Title   string   `json:"title"`
	Summary string   `json:"summary,omitempty"`
	URL     string   `json:"url,omitempty"`
	UsedBy  string   `json:"usedBy,omitempty"`
	Modules []string `json:"modules,omitempty"`
	Fn      CheckFn  `json:"-"`
	Solve   SolveFn  `json:"-"`
}

func (c *Check) Check(ctx context.Context, logger util.Logger) *Result {
	_, span, logger := telemetry.StartSpan(ctx, "doctor:check:"+c.Key, logger)
	defer span.Complete()

	r := NewResult(c, c.Key, c.Title, c.Summary)
	timer := util.TimerStart()
	r = c.Fn(ctx, r, logger)
	r.Duration = timer.End()
	r = c.Solve(ctx, r, logger)
	return r
}

type Checks []*Check

func (c Checks) Get(key string) *Check {
	return lo.FindOrElse(c, nil, func(x *Check) bool {
		return x.Key == key
	})
}
