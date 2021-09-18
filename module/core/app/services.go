// Package app $PF_IGNORE$
package app

import (
	"context"{{{ if .HasModule "database" }}}

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/database"
	"{{{ .Package }}}/app/database/migrate"
	"{{{ .Package }}}/app/log"{{{ end }}}
)

type Services struct {
  // add your stuff here
}

func NewServices(ctx context.Context, st *State) (*Services, error) {
	return &Services{}, nil
}
{{{ if .HasModule "database" }}}
func RunMigrations(context.Context, *zap.SugaredLogger) (*database.Service, error) {
	logger, _ := log.InitLogging(false)

	db, err := database.OpenDefaultPostgres(logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open database")
	}

	err = migrate.Migrate(context.Background(), db, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run database migrations")
	}

	return db, nil
}{{{ end }}}
