package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"projectforge.dev/projectforge/app/util"
)

func upF(ctx context.Context) error {
	if err := upgradeF(ctx); err != nil {
		return err
	}
	return updateF(ctx)
}

func upCmd() *cobra.Command {
	f := func(_ *cobra.Command, _ []string) error { return upF(rootCtx) }
	ret := &cobra.Command{Use: "up", Short: "Full update and upgrade of " + util.AppName, RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
