// Package controller - Content managed by Project Forge, see [projectforge.md] for details.
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
	ps := cutil.LoadPageState(as, w, r, key, _currentAppRootLogger)
	if err := initAppRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, w, r, f)
}

func ActSite(key string, w http.ResponseWriter, r *http.Request, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentSiteState
	ps := cutil.LoadPageState(as, w, r, key, _currentSiteRootLogger)
	ps.Menu = site.Menu(ps.Context, as, ps.Profile, ps.Logger)
	if err := initSiteRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, w, r, f)
}

func actComplete(key string, as *app.State, ps *cutil.PageState, w http.ResponseWriter, r *http.Request, f ActFn) {
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
	logger = logger.With("path", r.URL.Path, "method", ps.Method, "status", status)
	ps.Context = ctx

	if ps.ForceRedirect == "" || ps.ForceRedirect == r.URL.Path {
		redir, err = safeRun(f, as, ps)
		if err != nil {
			redir, err = handleError(key, as, ps, w, r, err)
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
	defer ps.Close()
	w.Header().Set("Server-Timing", fmt.Sprintf("server:dur=%.3f", elapsedMillis))
	logger = logger.With("elapsed", elapsedMillis)
	logger.Debugf("processed request in [%.3fms] (render: %.3fms)", elapsedMillis, ps.RenderElapsed)
}

func safeRun(f func(as *app.State, ps *cutil.PageState) (string, error), as *app.State, ps *cutil.PageState) (s string, e error) {
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
