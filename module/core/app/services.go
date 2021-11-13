// Package app $PF_IGNORE$
package app

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/database/migrate"
	"{{{ .Package }}}/app/queries/migrations"{{{ end }}}
)

type Services struct {
  // add your stuff here
}

func NewServices(ctx context.Context, st *State) (*Services, error) {
	{{{ if .HasModule "migration" }}}
	migrations.LoadMigrations()
	err := migrate.Migrate(ctx, st.DB, st.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}return &Services{}, nil
}
