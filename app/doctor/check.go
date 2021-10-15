package doctor

import (
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

type (
	checkFn func(r *Result, logger *zap.SugaredLogger) *Result
	solveFn func(r *Result, logger *zap.SugaredLogger) *Result
)

type Check struct {
	Key     string   `json:"key"`
	Section string   `json:"section"`
	Title   string   `json:"title"`
	Summary string   `json:"summary,omitempty"`
	URL     string   `json:"url,omitempty"`
	UsedBy  string   `json:"usedBy,omitempty"`
	Modules []string `json:"modules,omitempty"`
	Fn      checkFn  `json:"-"`
	Solve   solveFn  `json:"-"`
}

func (c *Check) Check(logger *zap.SugaredLogger) *Result {
	r := NewResult(c, c.Key, c.Title, c.Summary)
	start := util.TimerStart()
	r = c.Fn(r, logger)
	r.Duration = util.TimerEnd(start)
	r = c.Solve(r, logger)
	return r
}

func NewCheck(key string, section string, title string, summary string, mods ...string) *Check {
	return &Check{Key: key, Section: section, Title: title, Summary: summary, Modules: mods}
}

type Checks = []*Check
