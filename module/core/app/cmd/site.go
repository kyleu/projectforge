package cmd

import (
	"fmt"
	"runtime"

	"github.com/fasthttp/router"
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller"
	"$PF_PACKAGE$/app/filesystem"
	"$PF_PACKAGE$/app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const keySite = "site"

func siteCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the marketing site on port %d (by default)", util.AppPort)
	f := func(*cobra.Command, []string) error { return startSite(_flags.Port) }
	ret := &cobra.Command{Use: keySite, Short: short, RunE: f}
	return ret
}

func startSite(port uint16) error {
	err := initIfNeeded()
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	r, _, err := loadSite(_flags, port, _logger)
	if err != nil {
		return err
	}

	_, err = listenandserve(util.AppName, _flags.Address, port, r)
	return err
}

func loadSite(flags *Flags, port uint16, logger *zap.SugaredLogger) (*router.Router, *zap.SugaredLogger, error) {
	r := controller.SiteRoutes()

	f := filesystem.NewFileSystem(flags.ConfigDir, logger)

	st, err := app.NewState(flags.Debug, _buildInfo, r, f, logger)
	if err != nil {
		return nil, logger, err
	}

	controller.SetSiteState(st, logger)

	logger.Infof("started marketing site using address [%s:%d] on %s:%s", flags.Address, port, runtime.GOOS, runtime.GOARCH)

	return r, logger, nil
}
