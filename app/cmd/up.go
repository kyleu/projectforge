package cmd

import (
	"context"

	"github.com/muesli/coral"

	"projectforge.dev/projectforge/app/util"
)

func upF(ctx context.Context) error {
	if err := upgradeF(ctx); err != nil {
		return err
	}
	return updateF(ctx)
}

func upCmd() *coral.Command {
	f := func(_ *coral.Command, _ []string) error { return upF(context.Background()) }
	ret := &coral.Command{Use: "up", Short: "Full update and upgrade of " + util.AppName, RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
