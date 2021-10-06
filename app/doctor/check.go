package doctor

import (
	"github.com/kyleu/projectforge/app/util"
)

type (
	checkFn func(r *Result) *Result
	solveFn func(r *Result) (*Result, error)
)

type Check struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Summary string   `json:"summary,omitempty"`
	Modules []string `json:"modules,omitempty"`
	Fn      checkFn  `json:"-"`
	Solve   solveFn  `json:"-"`
}

func (c *Check) Check() *Result {
	r := NewResult(c, c.Key, c.Title, c.Summary)
	start := util.TimerStart()
	r = c.Fn(r)
	r.Duration = util.TimerEnd(start)
	return r
}

func NewCheck(key string, title string, summary string, mods ...string) *Check {
	return &Check{Key: key, Title: title, Summary: summary, Modules: mods}
}

type Checks = []*Check
