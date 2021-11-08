package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/log"
	"{{{ .Package }}}/app/upgrade"
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
	if _version == "" {
		_version = _buildInfo.Version
	}
	return upgrade.NewService(l).UpgradeIfNeeded(ctx, _buildInfo.Version, _version, _force)
}

func upgradeCmd() *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return upgradeF(context.Background(), args) }
	ret := &cobra.Command{Use: "upgrade", Short: "Upgrades " + util.AppKey + " to the latest published version", RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to update to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force update, even if same or earlier")
	return ret
}
