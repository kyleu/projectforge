package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/lib/upgrade"
	"{{{ .Package }}}/app/util"
)

var (
	_version = ""
	_force   = false
)

func upgradeF(ctx context.Context) error {
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	return upgrade.NewService(ctx, l).UpgradeIfNeeded(ctx, _buildInfo.Version, _version, _force)
}

func upgradeCmd() *cobra.Command {
	f := func(_ *cobra.Command, _ []string) error { return upgradeF(rootCtx) }
	ret := newCmd("upgrade", "Upgrades "+util.AppName+" to the latest published version", f)
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
