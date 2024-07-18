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
	util.RootLogger.Infof("updating " + util.AppName + " modules...")
	mSvc, err := module.NewService(ctx, _flags.ConfigDir, util.RootLogger)
	if err != nil {
		return err
	}
	for _, mod := range mSvc.Modules() {
		url := mod.URL
		var err error
		if url == "" {
			url, err = mSvc.AssetURL(ctx, mod.Key, util.RootLogger)
			if err != nil {
				return err
			}
		}
		err = mSvc.Download(ctx, mod.Key, url, util.RootLogger)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateCmd() *coral.Command {
	f := func(_ *coral.Command, _ []string) error { return updateF(context.Background()) }
	return &coral.Command{Use: "update", Short: "Refreshes downloaded assets such as modules", RunE: f}
}
