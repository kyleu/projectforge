package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/site"
	"projectforge.dev/projectforge/app/util"
)

type ActFn func(as *app.State, ps *cutil.PageState) (string, error)

func Act(key string, w http.ResponseWriter, r *http.Request, f ActFn) {
	as := _currentAppState
	wc := cutil.NewWriteCounter(w)
	ps := cutil.LoadPageState(as, wc, r, key, _currentAppRootLogger)
	if err := initAppRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(as, key, ps, wc, r, f)
}

func ActSite(key string, w http.ResponseWriter, r *http.Request, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentSiteState
	wc := cutil.NewWriteCounter(w)
	ps := cutil.LoadPageState(as, wc, r, key, _currentSiteRootLogger)
	ps.Menu = site.Menu(ps.Context, as, ps.Profile, ps.Logger)
	if err := initSiteRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(as, key, ps, ps.W, r, f)
}

func actComplete(as *app.State, key string, ps *cutil.PageState, w *cutil.WriteCounter, r *http.Request, f ActFn) {
	err := ps.Clean(r, as)
	if err != nil {
		ps.Logger.Warnf("error while cleaning request, somehow: %+v", err)
	}
	status := http.StatusOK
	cutil.WriteCORS(w)
	var redir string
	logger := ps.Logger
	ctx := ps.Context
	if !telemetry.SkipControllerMetrics {
		var span *telemetry.Span
		ctx, span, logger = telemetry.StartSpan(ps.Context, "controller."+key, ps.Logger)
		defer span.Complete()
	}
	logger = logger.With("path", ps.URI.Path, "method", ps.Method, "status", status)
	ps.Context = ctx

	if ps.ForceRedirect == "" || ps.ForceRedirect == ps.URI.Path {
		redir, err = safeRun(as, f, ps)
		if err != nil {
			redir, err = handleError(as, key, ps, r, err)
			if err != nil {
				ps.Logger.Warnf("unable to handle error: %+v", err)
			}
		}
	} else {
		redir = ps.ForceRedirect
	}
	if redir != "" {
		w.Header().Set("Location", redir)
		w.WriteHeader(http.StatusFound)
	}
	elapsedMillis := float64((util.TimeCurrentNanos()-ps.Started.UnixNano())/int64(time.Microsecond)) / float64(1000)
	ps.ResponseBytes = w.Count()
	defer ps.Close()
	w.Header().Set("Server-Timing", fmt.Sprintf("server:dur=%.3f", elapsedMillis))
	logger = logger.With("elapsed", elapsedMillis)
	if ps.Transport == "ws" {
		logger.Debugf("closed websocket request after [%.3fms]", elapsedMillis)
	} else {
		logger.Debugf("processed request in [%.3fms] (render: %.3fms, response: %s)", elapsedMillis, ps.RenderElapsed, util.ByteSizeSI(ps.ResponseBytes))
	}
}

func safeRun(as *app.State, f func(as *app.State, ps *cutil.PageState) (string, error), ps *cutil.PageState) (s string, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if recoverErr, ok := rec.(error); ok {
				e = errors.Wrap(recoverErr, "panic")
			} else {
				e = errors.Errorf("controller encountered panic recovery of type [%T]: %s", rec, fmt.Sprint(rec))
			}
		}
	}()
	s, e = f(as, ps)
	return
}
