package doctor

type checkFn func(r *Result) *Result

type Check struct {
	Key     string   `json:"key"`
	Title   string   `json:"title"`
	Summary string   `json:"summary,omitempty"`
	Modules []string `json:"modules,omitempty"`
	Fn      checkFn  `json:"-"`
}

func NewCheck(key string, title string, summary string, mods ...string) *Check {
	return &Check{Key: key, Title: title, Summary: summary, Modules: mods}
}

type Checks = []*Check
