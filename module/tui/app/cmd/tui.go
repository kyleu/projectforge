//go:build !js
// +build !js

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"{{{ .Package }}}/app/controller/tui"
	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/util"
)

const keyTUI = "tui"

func tuiCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the terminal UI (and the http server on port %d)", util.AppPort)
	f := func(*cobra.Command, []string) error { return runTUI(rootCtx, _flags) }
	ret := &cobra.Command{Use: keyTUI, Short: short, RunE: f}
	return ret
}

func runTUI(ctx context.Context, flags *Flags) error {
	var t *tui.TUI
	logFn := func(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) {
		if t != nil {
			t.AddLog(level, occurred, loggerName, message, caller, stack, fields)
		}
	}

	l, err := log.InitQuietLogging(zapcore.DebugLevel, logFn)
	if err != nil {
		return errors.Wrap(err, "error initializing logging")
	}
	util.RootLogger = l.Sugar()

	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	st, r, logger, err := loadServer(ctx, flags, logger)
	if err != nil {
		return err
	}

	var serverURL string
	var serverErr string
	var srv *http.Server
	srvErr := make(chan error, 1)

	port, ln, err := listen(ctx, flags.Address, flags.Port)
	if err != nil {
		serverErr = errors.Wrap(err, "unable to start http server").Error()
	} else {
		srv = newHTTPServer(r)
		go func() {
			serveErr := srv.Serve(ln)
			if serveErr == nil || errors.Is(serveErr, http.ErrServerClosed) {
				srvErr <- nil
				return
			}
			srvErr <- errors.Wrap(serveErr, "http server exited")
		}()
		serverURL = fmt.Sprintf("http://%s:%d", flags.Address, port)
		select {
		case serveErr := <-srvErr:
			if serveErr != nil {
				serverErr = serveErr.Error()
				serverURL = ""
			}
		default:
		}
	}

	t, err = tui.NewTUI(st, serverURL, serverErr, logger)
	if err != nil {
		return err
	}
	tuiErr := t.Run(ctx, logger)

	if srv != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
		}
		cancel()
	}

	select {
	case err := <-srvErr:
		if err != nil {
			if tuiErr != nil {
				return errors.Wrapf(tuiErr, "tui exited (also: %s)", err.Error())
			}
			return err
		}
	case <-time.After(250 * time.Millisecond):
	}

	err = st.Close(ctx, logger)
	if err != nil {
		return errors.Wrap(err, "unable to close application")
	}

	return tuiErr
}
