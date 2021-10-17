package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func versionF(context.Context, []string) error {
	println(_buildInfo.Version)
	return nil
}

func versionCmd() *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return versionF(context.Background(), args) }
	return &cobra.Command{Use: "version", Short: "Displays the version and exits", RunE: f}
}
