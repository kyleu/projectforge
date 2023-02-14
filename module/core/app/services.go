// Package app - $PF_GENERATE_ONCE$
package app

import (
	"context"{{{ if .HasModule "migration" }}}

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database/migrate"
	"{{{ .Package }}}/queries/migrations"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

type Services struct {
	{{{ if.HasModule "export" }}}
	// $PF_INJECT_START(services)$
	// $PF_INJECT_END(services)$
	{{{ end }}}// add your dependencies here
}

func NewServices(ctx context.Context, st *State, logger util.Logger) (*Services, error) {
	{{{ if .HasModule "migration" }}}migrations.LoadMigrations(st.Debug)
	err := migrate.Migrate(ctx, st.DB, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}
	{{{ end }}}return &Services{
	{{{ if.HasModule "export" }}}	// $PF_INJECT_START(refs)$
		// $PF_INJECT_END(refs)${{{ end }}}
	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
