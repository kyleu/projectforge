package sandbox

import (
	"github.com/kyleu/projectforge/app"
	"go.uber.org/zap"
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

var AllSandboxes = Sandboxes{example, testbed}

var example = &Sandbox{Key: "example", Title: "Example", Icon: "play", Run: func(st *app.State, logger *zap.SugaredLogger) (interface{}, error) {
	return "a work in progress...", nil
}}
