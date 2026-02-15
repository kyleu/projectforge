//go:build !js
// +build !js

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	"projectforge.dev/projectforge/app/controller/tui"
	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/util"
)

const keyTUI = "tui"

func tuiCmd() *coral.Command {
	short := fmt.Sprintf("Starts the terminal UI (and the http server on port %d)", util.AppPort)
	f := func(*coral.Command, []string) error { return runTUI(rootCtx, _flags) }
	ret := &coral.Command{Use: keyTUI, Short: short, RunE: f}
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

	port, ln, err := listen(ctx, flags.Address, flags.Port)
	if err != nil {
		return err
	}
	srv := newHTTPServer(r)
	srvErr := make(chan error, 1)
	go func() {
		err := srv.Serve(ln)
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			srvErr <- nil
			return
		}
		srvErr <- errors.Wrap(err, "http server exited")
	}()

	serverURL := fmt.Sprintf("http://%s:%d", flags.Address, port)
	t, err = tui.NewTUI(st, serverURL, logger)
	if err != nil {
		return err
	}
	tuiErr := t.Run(ctx, logger)
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		_ = srv.Close()
	}
	cancel()

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
