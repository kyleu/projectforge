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

func upgradeF(ctx context.Context, _ []string) error {
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	return upgrade.NewService(ctx, l).UpgradeIfNeeded(ctx, _buildInfo.Version, _version, _force)
}

func upgradeCmd() *coral.Command {
	f := func(cmd *coral.Command, args []string) error { return upgradeF(context.Background(), args) }
	ret := &coral.Command{Use: "upgrade", Short: "Upgrades " + util.AppKey + " to the latest published version", RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to update to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force update, even if same or earlier")
	return ret
}
