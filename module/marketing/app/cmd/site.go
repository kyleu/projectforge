package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/routes"
	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

const keySite = "site"

func siteCmd() *coral.Command {
	short := fmt.Sprintf("Starts the marketing site on port %d (by default)", util.AppPort)
	f := func(*coral.Command, []string) error { return startSite(rootCtx, _flags) }
	ret := &coral.Command{Use: keySite, Short: short, RunE: f}
	return ret
}

func startSite(ctx context.Context, flags *Flags) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	r, _, err := loadSite(ctx, flags, logger)
	if err != nil {
		return err
	}

	_, err = listenAndServe(flags.Address, flags.Port, r)
	return err
}

func loadSite(ctx context.Context, flags *Flags, logger util.Logger) (http.Handler, util.Logger, error) {
	r, err := routes.SiteRoutes(logger)
	if err != nil {
		return nil, logger, err
	}
	f, err := filesystem.NewFileSystem(flags.ConfigDir, false, "")
	if err != nil {
		return nil, logger, err
	}

	telemetryDisabled := util.GetEnvBool("disable_telemetry", false)
	st, err := app.NewState(ctx, flags.Debug, _buildInfo, f, !telemetryDisabled, flags.Port, logger)
	if err != nil {
		return nil, logger, err
	}

	err = controller.SetSiteState(ctx, st, logger)
	if err != nil {
		return nil, logger, err
	}
	logger.Infof("started marketing site using address [http://%s:%d] on %s:%s", flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)
	return r, logger, nil
}
