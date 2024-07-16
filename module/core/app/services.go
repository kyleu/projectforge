// Package app - $PF_GENERATE_ONCE$
package app // import {{{ .Package }}}

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"{{{ end }}}
	{{{ if .HasModule "audit" }}}
	"{{{ .Package }}}/app/lib/audit"{{{ end }}}{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database/migrate"{{{ end }}}
	"{{{ .Package }}}/app/util"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/queries/migrations"{{{ end }}}
)

type Services struct {
	CoreServices{{{ if .HasExport }}}
	GeneratedServices{{{ end }}}

	// add your dependencies here
}

func NewServices(ctx context.Context, st *State, logger util.Logger) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}{{{ if .HasModule "audit" }}}auditSvc := audit.NewService(st.DB, logger)
	{{{ end }}}return &Services{
		CoreServices:      initCoreServices(ctx, st{{{ if .HasModule "audit" }}}, auditSvc{{{ end }}}, logger),{{{ if .HasExport }}}
		GeneratedServices: initGeneratedServices(ctx, st{{{ if .HasModule "audit" }}}, auditSvc{{{ end }}}, logger),{{{ end }}}
	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
