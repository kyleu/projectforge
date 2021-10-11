package doctor

import (
	"github.com/kyleu/projectforge/app/util"
)

type (
	checkFn func(r *Result) *Result
	solveFn func(r *Result) *Result
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

func (c *Check) Check() *Result {
	r := NewResult(c, c.Key, c.Title, c.Summary)
	start := util.TimerStart()
	r = c.Fn(r)
	r.Duration = util.TimerEnd(start)
	r = c.Solve(r)
	return r
}

func NewCheck(key string, section string, title string, summary string, mods ...string) *Check {
	return &Check{Key: key, Section: section, Title: title, Summary: summary, Modules: mods}
}

type Checks = []*Check
