package cmd

import (
	"context"

	"github.com/muesli/coral"

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

func upgradeCmd() *coral.Command {
	f := func(cmd *coral.Command, _ []string) error { return upgradeF(context.Background()) }
	ret := &coral.Command{Use: "upgrade", Short: "Upgrades " + util.AppName + " to the latest published version", RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
