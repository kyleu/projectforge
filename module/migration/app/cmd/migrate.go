package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/lib/database/migrate"
	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/queries/migrations"
)

func migrateCmd() *cobra.Command {
	f := func(*cobra.Command, []string) error { return runMigrations() }
	ret := &cobra.Command{Use: "migrate", Short: "Runs database migrations and exits", RunE: f}
	return ret
}

func runMigrations() error {
	logger, _ := log.InitLogging(false)
	db, err := database.OpenDefaultPostgres(context.Background(), logger)
	if err != nil {
		return errors.Wrap(err, "unable to open database")
	}
	migrations.LoadMigrations(_flags.Debug)
	err = migrate.Migrate(context.Background(), db, logger)
	if err != nil {
		return errors.Wrap(err, "unable to run database migrations")
	}
	return nil
}
