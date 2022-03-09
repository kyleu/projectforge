// Content managed by Project Forge, see [projectforge.md] for details.
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

func upgradeCmd() *coral.Command {
	f := func(cmd *coral.Command, args []string) error { return upgradeF(context.Background(), args) }
	ret := &coral.Command{Use: "upgrade", Short: "Upgrades " + util.AppKey + " to the latest published version", RunE: f}
	ret.PersistentFlags().StringVar(&_version, "version", "", "version number to update to")
	ret.PersistentFlags().BoolVarP(&_force, "force", "f", false, "force update, even if same or earlier")
	return ret
}
