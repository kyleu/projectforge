package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/routes"
	"projectforge.dev/projectforge/app/util"
)

const keyServer = "server"

func serverCmd() *coral.Command {
	short := fmt.Sprintf("Starts the http server on port %d (by default)", util.AppPort)
	f := func(*coral.Command, []string) error { return startServer(_flags) }
	ret := &coral.Command{Use: keyServer, Short: short, RunE: f}
	return ret
}

func startServer(flags *Flags) error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	st, r, logger, err := loadServer(flags, util.RootLogger)
	if err != nil {
		return err
	}

	_, err = listenandserve(flags.Address, flags.Port, r)
	if err != nil {
		return err
	}

	err = st.Close(context.Background(), logger)
	if err != nil {
		return errors.Wrap(err, "unable to close application")
	}

	return nil
}

func loadServer(flags *Flags, logger util.Logger) (*app.State, http.Handler, util.Logger, error) {
	st, err := app.Bootstrap(_buildInfo, _flags.ConfigDir, _flags.Port, _flags.Debug, logger)
	if err != nil {
		return nil, nil, logger, err
	}
	logger.Infof("started %s v%s using address [http://%s:%d] on %s:%s", util.AppName, _buildInfo.Version, flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)

	err = controller.SetAppState(st, logger)
	if err != nil {
		return nil, nil, logger, err
	}
	r, err := routes.AppRoutes(st, logger)
	if err != nil {
		return nil, nil, logger, err
	}

	return st, r, logger, nil
}
