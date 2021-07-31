package sandbox

import (
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
)

type runFn func(st *app.State, logger *zap.SugaredLogger) (interface{}, error)

type Sandbox struct {
	Key   string `json:"key,omitempty"`
	Title string `json:"title,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Run   runFn  `json:"-"`
}

type Sandboxes []*Sandbox

func (s Sandboxes) Get(key string) *Sandbox {
	for _, v := range s {
		if v.Key == key {
			return v
		}
	}
	return nil
}

// $PF_SECTION_START(sandboxes)$
var AllSandboxes = Sandboxes{testbed}

// $PF_SECTION_END(sandboxes)$
