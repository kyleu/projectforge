package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/log"
)

func migrateCmd() *cobra.Command {
	f := func(*cobra.Command, []string) error { return runMigrations() }
	ret := &cobra.Command{Use: "migrate", Short: "Runs database migrations and exits", RunE: f}
	return ret
}

func runMigrations() error {
	logger, _ := log.InitLogging(false)
	_, err := app.RunMigrations(context.Background(), logger)
	return err
}
