package cmd

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/filesystem"
	"{{{ .Package }}}/app/util"
)

const keySite = "site"

func siteCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the marketing site on port %d (by default)", util.AppPort)
	f := func(*cobra.Command, []string) error { return startSite(_flags) }
	ret := &cobra.Command{Use: keySite, Short: short, RunE: f}
	return ret
}

func startSite(flags *Flags) error {
	err := initIfNeeded()
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	r, _, err := loadSite(flags, _logger)
	if err != nil {
		return err
	}

	_, err = listenandserve(util.AppName, flags.Address, flags.Port, r)
	return err
}

func loadSite(flags *Flags, logger *zap.SugaredLogger) (fasthttp.RequestHandler, *zap.SugaredLogger, error) {
	m, r := controller.SiteRoutes()
	f := filesystem.NewFileSystem(flags.ConfigDir, logger)

	st, err := app.NewState(flags.Debug, _buildInfo, f, m, logger)
	if err != nil {
		return nil, logger, err
	}

	controller.SetSiteState(st)
	logger.Infof("started marketing site using address [%s:%d] on %s:%s", flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)
	return r, logger, nil
}