package mvc

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/theme"
	"{{{ .Package }}}/app/util"
)

type State struct {
	App       *app.State
	Context   context.Context
	Logger    util.Logger
	ServerURL string
	ServerErr string
	Theme     *theme.Theme
	LogTail   func(limit int) []string
}

func NewState(ctx context.Context, st *app.State, serverURL string, serverErr string, logger util.Logger, logTail func(limit int) []string) *State {
	if ctx == nil {
		ctx = context.Background()
	}
	ret := &State{App: st, Context: ctx, ServerURL: serverURL, ServerErr: serverErr, Logger: logger, Theme: theme.Default, LogTail: logTail}
	if st != nil && st.Themes != nil {
		ret.Theme = st.Themes.Get("default", logger)
	}
	return ret
}
