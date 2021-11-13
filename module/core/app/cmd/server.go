package cmd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/database"{{{ end }}}
	"{{{ .Package }}}/app/filesystem"
	"{{{ .Package }}}/app/util"
)

const keyServer = "server"

func serverCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the http server on port %d (by default)", util.AppPort)
	f := func(*cobra.Command, []string) error { return startServer(_flags) }
	ret := &cobra.Command{Use: keyServer, Short: short, RunE: f}
	return ret
}

func startServer(flags *Flags) error {
	err := initIfNeeded()
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	r, _, err := loadServer(flags, _logger)
	if err != nil {
		return err
	}

	_, err = listenandserve(util.AppName, flags.Address, flags.Port, r)
	return err
}

func loadServer(flags *Flags, logger *zap.SugaredLogger) (fasthttp.RequestHandler, *zap.SugaredLogger, error) {
	r := controller.AppRoutes()
	f := filesystem.NewFileSystem(flags.ConfigDir, logger)
	st, err := app.NewState(flags.Debug, _buildInfo, f, logger)
	if err != nil {
		return nil, logger, err
	}{{{ if .HasModule "migration" }}}

	db, err := database.OpenDefaultPostgres(logger)
	if err != nil {
		return nil, logger, errors.Wrap(err, "unable to open database")
	}
	st.DB = db{{{ end }}}

	svcs, err := app.NewServices(context.Background(), st)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating services")
	}
	st.Services = svcs

	controller.SetAppState(st)

	logger.Infof("started %s v%s using address [%s:%d] on %s:%s", util.AppName, _buildInfo.Version, flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)

	return r, logger, nil
}
