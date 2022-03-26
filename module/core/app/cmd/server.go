package cmd

import (
	"context"
	"fmt"
	"runtime"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"{{{ if .HasModule "migration" }}}
	"{{{ .Package }}}/app/lib/database"{{{ end }}}
	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
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
	telemetryEnabled := util.GetEnvBool("disable_telemetry", false)
	st, err := app.NewState(flags.Debug, _buildInfo, f, telemetryEnabled, logger)
	if err != nil {
		return nil, logger, err
	}

	ctx, span, logger := telemetry.StartSpan(context.Background(), "app:init", logger)
	defer span.Complete(){{{ if .HasModule "migration" }}}{{{ if .HasModule "postgres" }}}

	db, err := database.OpenDefaultPostgres(ctx, logger){{{ else }}}{{{ if .HasModule "sqlite" }}}
	db, err := database.OpenDefaultSQLite(ctx, logger){{{ end }}}{{{ end }}}
	if err != nil {
		return nil, logger, errors.Wrap(err, "unable to open database")
	}
	st.DB = db{{{ end }}}

	svcs, err := app.NewServices(ctx, st)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating services")
	}
	st.Services = svcs

	controller.SetAppState(st)

	logger.Infof("started %s v%s using address [%s:%d] on %s:%s", util.AppName, _buildInfo.Version, flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)

	return r, logger, nil
}
