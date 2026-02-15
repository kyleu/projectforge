package mvc

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

type State struct {
	App       *app.State
	Context   context.Context
	Logger    util.Logger
	ServerURL string
	Theme     *theme.Theme
}

func NewState(ctx context.Context, st *app.State, serverURL string, logger util.Logger) *State {
	if ctx == nil {
		ctx = context.Background()
	}
	ret := &State{App: st, Context: ctx, ServerURL: serverURL, Logger: logger, Theme: theme.Default}
	if st != nil && st.Themes != nil {
		ret.Theme = st.Themes.Get("default", logger)
	}
	return ret
}
