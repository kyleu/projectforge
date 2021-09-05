package sandbox

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/menu"
)

type runFn func(ctx context.Context, st *app.State, logger *zap.SugaredLogger) (interface{}, error)

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

func Menu() *menu.Item {
	ret := make(menu.Items, 0, len(AllSandboxes))
	for _, s := range AllSandboxes {
		desc := fmt.Sprintf("Sandbox [%s]", s.Key)
		rt := fmt.Sprintf("/admin/sandbox/%s", s.Key)
		ret = append(ret, &menu.Item{Key: s.Key, Title: s.Title, Icon: s.Icon, Description: desc, Route: rt})
	}
	desc := "Playgrounds for testing new features"
	return &menu.Item{Key: "sandbox", Title: "Sandboxes", Description: desc, Icon: "play", Route: "/admin/sandbox", Children: ret}
}
