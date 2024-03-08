package cmd

import (
	"context"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/database/migrate"
	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/queries/migrations"
)

func migrateCmd() *coral.Command {
	f := func(*coral.Command, []string) error { return runMigrations(context.Background()) }
	ret := &coral.Command{Use: "migrate", Short: "Runs database migrations and exits", RunE: f}
	return ret
}

func runMigrations(ctx context.Context) error {
	logger, _ := log.InitLogging(false){{{ if .PostgreSQL }}}
	db, err := database.OpenDefaultPostgres(ctx, logger){{{ else }}}{{{ if .SQLite }}}
	db, err := database.OpenDefaultSQLite(ctx, logger){{{ else }}}{{{ if .SQLServer }}}
	db, err := database.OpenDefaultSQLServer(ctx, logger){{{ end }}}{{{ end }}}{{{ end }}}
	if err != nil {
		return errors.Wrap(err, "unable to open database")
	}
	migrations.LoadMigrations(_flags.Debug)
	err = migrate.Migrate(ctx, db, logger)
	if err != nil {
		return errors.Wrap(err, "unable to run database migrations")
	}
	return nil
}
