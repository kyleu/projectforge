// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"context"

	"github.com/muesli/coral"

	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/lib/upgrade"
	"projectforge.dev/projectforge/app/util"
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
	f := func(_ *coral.Command, _ []string) error { return upgradeF(context.Background()) }
	ret := &coral.Command{Use: "upgrade", Short: "Upgrades " + util.AppName + " to the latest published version", RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to upgrade to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force upgrade, even if same or earlier")
	return ret
}
