// Package app - $PF_IGNORE$
package app

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"{{{ end }}}

	"{{{ .Package }}}/app/lib/database/migrate"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/queries/migrations"{{{ end }}}
)

type Services struct {
	{{{ if.HasModule "export" }}}
	// $PF_INJECT_START(services)$
	// $PF_INJECT_END(services)$
	{{{ end }}}// add your dependencies here
}

func NewServices(ctx context.Context, st *State) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, st.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}return &Services{
	{{{ if.HasModule "export" }}}	// $PF_INJECT_START(refs)$
		// $PF_INJECT_END(refs)${{{ end }}}
	}, nil
}
