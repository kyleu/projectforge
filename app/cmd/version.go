// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"context"

	"github.com/muesli/coral"
)

func versionF(_ context.Context, _ []string) error {
	println(_buildInfo.Version) // nolint
	return nil
}

func versionCmd() *coral.Command {
	f := func(cmd *coral.Command, args []string) error { return versionF(context.Background(), args) }
	return &coral.Command{Use: "version", Short: "Displays the version and exits", RunE: f}
}
