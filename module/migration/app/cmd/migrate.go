package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/database"
	"{{{ .Package }}}/app/database/migrate"
	"{{{ .Package }}}/app/log"
	"{{{ .Package }}}/queries/migrations"
)

func migrateCmd() *cobra.Command {
	f := func(*cobra.Command, []string) error { return runMigrations() }
	ret := &cobra.Command{Use: "migrate", Short: "Runs database migrations and exits", RunE: f}
	return ret
}

func runMigrations() error {
	logger, _ := log.InitLogging(false)
	db, err := database.OpenDefaultPostgres(logger)
	if err != nil {
		return errors.Wrap(err, "unable to open database")
	}
	migrations.LoadMigrations()
	err = migrate.Migrate(context.Background(), db, logger)
	if err != nil {
		return errors.Wrap(err, "unable to run database migrations")
	}
	return nil
}
