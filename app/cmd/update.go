package cmd

import (
	"context"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func updateF(ctx context.Context, args []string) error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	_logger.Infof("updating " + util.AppName + " modules...")
	fs := filesystem.NewFileSystem(_flags.ConfigDir, _logger)
	mSvc := module.NewService(fs, _logger)
	for _, mod := range mSvc.Modules() {
		err := mSvc.Download(mod.Key, mod.URL)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateCmd() *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return updateF(context.Background(), args) }
	return &cobra.Command{Use: "update", Short: "Refreshes downloaded assets such as modules", RunE: f}
}
