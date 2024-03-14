// Package app - $PF_GENERATE_ONCE$
package app // import {{{ .Package }}}

import (
	"context"

	{{{ .ServicesImports }}}
)

{{{ .ServicesDefinition }}}

func NewServices(ctx context.Context, st *State, logger util.Logger) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}{{{ if .HasModule "audit" }}}auditSvc := audit.NewService(st.DB, logger)
	{{{ end }}}return {{{ .ServicesConstructor }}}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
