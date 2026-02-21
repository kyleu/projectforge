package cmd

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/routes"
	"projectforge.dev/projectforge/app/util"
)

const keyServer = "server"

func serverCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the http server on port %d (by default)", util.AppPort)
	f := func(*cobra.Command, []string) error { return startServer(rootCtx, _flags) }
	ret := newCmd(keyServer, short, f)
	return ret
}

func startServer(ctx context.Context, flags *Flags) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	st, r, logger, err := loadServer(ctx, flags, logger)
	if err != nil {
		return err
	}

	// $PF_SECTION_START(extraserver)$
	// $PF_SECTION_END(extraserver)$

	_, err = listenAndServe(ctx, flags.Address, flags.Port, r)
	if err != nil {
		return err
	}

	err = st.Close(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "unable to close application")
	}

	return nil
}

func loadServer(ctx context.Context, flags *Flags, logger util.Logger) (*app.State, http.Handler, util.Logger, error) {
	st, err := app.Bootstrap(ctx, _buildInfo, _flags.ConfigDir, _flags.Port, _flags.Debug, logger)
	if err != nil {
		return nil, nil, logger, err
	}
	logger.Infof("started %s v%s using address [http://%s:%d] on %s:%s", util.AppName, _buildInfo.Version, flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)

	err = controller.SetAppState(ctx, st, logger)
	if err != nil {
		return nil, nil, logger, err
	}
	r, err := routes.AppRoutes(st, logger)
	if err != nil {
		return nil, nil, logger, err
	}

	return st, r, logger, nil
}
