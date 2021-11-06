package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/log"
	"{{{ .Package }}}/app/upgrade"
	"{{{ .Package }}}/app/util"
)

func upgradeF(ctx context.Context, cmds []string) error {
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	force := false
	for _, x := range cmds {
		if x == "force" {
			force = true
		}
	}
	return upgrade.NewService(l).UpgradeIfNeeded(ctx, _buildInfo.Version, force)
}

func upgradeCmd() *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return upgradeF(context.Background(), args) }
	return &cobra.Command{Use: "upgrade", Short: "Upgrades " + util.AppKey + " to the latest published version", RunE: f}
}
