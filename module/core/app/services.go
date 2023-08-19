// Package app - $PF_GENERATE_ONCE$
package app

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database/migrate"{{{ end}}}{{{ if.HasModule "process" }}}
	"{{{ .Package }}}/app/lib/exec"{{{ end }}}{{{ if.HasModule "scripting" }}}
	"{{{ .Package }}}/app/lib/scripting"{{{ end }}}{{{ if.HasModule "websocket" }}}
	"{{{ .Package }}}/app/lib/websocket"{{{ end }}}{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/queries/migrations"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

type Services struct {
	{{{ if.HasModule "export" }}}// $PF_INJECT_START(services)$
	// $PF_INJECT_END(services)${{{ end }}}{{{ if.HasModule "process" }}}
	Exec   *exec.Service{{{ end }}}{{{ if.HasModule "scripting" }}}
	Script *scripting.Service{{{ end }}}{{{ if.HasModule "websocket" }}}
	Socket *websocket.Service{{{ end }}}
	// add your dependencies here
}

func NewServices(ctx context.Context, st *State, logger util.Logger) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}return &Services{
		{{{ if.HasModule "process" }}}Exec:   exec.NewService(),
		{{{ end }}}{{{ if.HasModule "scripting" }}}Script: scripting.NewService(st.FS, "scripts"),
		{{{ end }}}{{{ if.HasModule "websocket" }}}Socket: websocket.NewService(nil, nil, nil),
		{{{ end}}}{{{ if.HasModule "export" }}}// $PF_INJECT_START(refs)$
		// $PF_INJECT_END(refs)${{{ end }}}
	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
