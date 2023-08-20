package cmd

import (
	"context"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func updateF(ctx context.Context) error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	_logger.Infof("updating " + util.AppName + " modules...")
	mSvc, err := module.NewService(ctx, _flags.ConfigDir, _logger)
	if err != nil {
		return err
	}
	for _, mod := range mSvc.Modules() {
		url := mod.URL
		var err error
		if url == "" {
			url, err = mSvc.AssetURL(ctx, mod.Key, _logger)
			if err != nil {
				return err
			}
		}
		err = mSvc.Download(ctx, mod.Key, url, _logger)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateCmd() *coral.Command {
	f := func(cmd *coral.Command, _ []string) error { return updateF(context.Background()) }
	return &coral.Command{Use: "update", Short: "Refreshes downloaded assets such as modules", RunE: f}
}
