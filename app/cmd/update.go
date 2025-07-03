package cmd

import (
	"context"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func updateF(ctx context.Context) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	logger.Infof("updating " + util.AppName + " modules...")
	mSvc, err := module.NewService(ctx, _flags.ConfigDir, logger)
	if err != nil {
		return err
	}
	keys := util.ArraySorted(mSvc.Modules().Keys())
	for _, key := range keys {
		mod := mSvc.Modules().Get(key)
		url := mod.URL
		var err error
		if url == "" {
			url, err = mSvc.AssetURL(ctx, mod.Key, logger)
			if err != nil {
				return err
			}
		}
		err = mSvc.Download(ctx, mod.Key, url, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateCmd() *coral.Command {
	f := func(_ *coral.Command, _ []string) error { return updateF(rootCtx) }
	return &coral.Command{Use: "update", Short: "Refreshes downloaded assets such as modules", RunE: f}
}
