// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"context"

	"github.com/muesli/coral"

	"projectforge.dev/projectforge/app/util"
)

func up(ctx context.Context) error {
	if err := upgradeF(ctx); err != nil {
		return err
	}
	return updateF(ctx)
}

func upCmd() *coral.Command {
	f := func(cmd *coral.Command, _ []string) error { return upgradeF(context.Background()) }
	ret := &coral.Command{Use: "up", Short: "Full update and upgrade of " + util.AppName, RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
