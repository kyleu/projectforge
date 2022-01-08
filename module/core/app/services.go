// Package app $PF_IGNORE$
package app

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/queries/migrations"{{{ end }}}
	"{{{ .Package }}}/app/lib/database/migrate"
)

type Services struct {
  // add your stuff here
}

func NewServices(ctx context.Context, st *State) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, st.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}return &Services{}, nil
}
