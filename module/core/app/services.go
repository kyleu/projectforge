// Package app - $PF_GENERATE_ONCE$
package app // import {{{ .Package }}}

import (
	"context"

	{{{ if .HasModule "audit" }}}"{{{ .Package }}}/app/lib/audit"{{{ end}}}{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database/migrate"{{{ end}}}{{{ if.HasModule "process" }}}
	"{{{ .Package }}}/app/lib/exec"{{{ end }}}{{{ if .HasModule "help" }}}
	"{{{ .Package }}}/app/lib/help"{{{ end }}}{{{ if.HasModule "scripting" }}}
	"{{{ .Package }}}/app/lib/scripting"{{{ end }}}{{{ if.HasModule "websocket" }}}
	"{{{ .Package }}}/app/lib/websocket"{{{ end }}}{{{ if.HasUser }}}
	"{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/queries/migrations"{{{ end }}}
)

type Services struct {
	{{{ if.HasModule "export" }}}GeneratedServices

{{{ end }}}{{{ if.HasModule "audit" }}}	Audit  *audit.Service{{{ end }}}{{{ if.HasModule "process" }}}
	Exec   *exec.Service{{{ end }}}{{{ if.HasModule "scripting" }}}
	Script *scripting.Service{{{ end }}}{{{ if.HasModule "websocket" }}}
	Socket *websocket.Service{{{ end }}}{{{ if.HasUser }}}
	User   *user.Service{{{ end }}}{{{ if.HasModule "help" }}}
	Help   *help.Service{{{ end }}}
	// add your dependencies here
}

func NewServices(ctx context.Context, st *State, logger util.Logger) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, logger)
	if err != nil {
		return nil, err
	}
	{{{ end }}}{{{ if .HasModule "audit" }}}auditSvc := audit.NewService(st.DB, logger)
	{{{ end }}}return &Services{
		{{{ if.HasModule "export" }}}GeneratedServices: initGeneratedServices(ctx, st.DB{{{ if .HasModule "audit" }}}, auditSvc{{{ end }}}, logger),

		{{{ end }}}{{{ if .HasModule "audit" }}}Audit:  auditSvc,
		{{{ end }}}{{{ if.HasModule "process" }}}Exec:   exec.NewService(),
		{{{ end }}}{{{ if.HasModule "scripting" }}}Script: scripting.NewService(st.Files, "scripts"),
		{{{ end }}}{{{ if.HasModule "websocket" }}}Socket: websocket.NewService(nil, nil, nil),
		{{{ end }}}{{{ if.HasUser }}}User:   user.NewService(st.Files, logger),
		{{{ end}}}{{{ if.HasModule "help" }}}Help:   help.NewService(logger),
{{{ end}}}	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
