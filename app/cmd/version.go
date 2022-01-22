// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func versionF(_ context.Context, _ []string) error {
	println(_buildInfo.Version) // nolint
	return nil
}

func versionCmd() *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return versionF(context.Background(), args) }
	return &cobra.Command{Use: "version", Short: "Displays the version and exits", RunE: f}
}
