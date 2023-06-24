package controller

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/telemetry"{{{ if .HasModule "oauth" }}}
	"{{{ .Package }}}/app/lib/user"{{{ end }}}{{{ if.HasModule "marketing" }}}
	"{{{ .Package }}}/app/site"{{{ end }}}
)

func Act(key string, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentAppState
	ps := cutil.LoadPageState(as, rc, key, _currentAppRootLogger){{{ if .HasModule "oauth" }}}
	if allowed, reason := user.Check(string(ps.URI.Path()), ps.Accounts); !allowed {
		f = Unauthorized(rc, reason, ps.Accounts)
	}{{{ end }}}
	if err := initAppRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, rc, f)
}
{{{ if.HasModule "marketing" }}}
func ActSite(key string, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentSiteState
	ps := cutil.LoadPageState(as, rc, key, _currentSiteRootLogger)
	ps.Menu = site.Menu(ps.Context, as, ps.Profile{{{ if .HasModule "oauth" }}}, ps.Accounts{{{ end }}}, ps.Logger){{{ if .HasModule "oauth" }}}
	if allowed, reason := user.Check(string(ps.URI.Path()), ps.Accounts); !allowed {
		f = Unauthorized(rc, reason, ps.Accounts)
	}{{{ end }}}
	if err := initSiteRequest(as, ps); err != nil {
		ps.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, rc, f)
}
{{{ end }}}
func actComplete(key string, as *app.State, ps *cutil.PageState, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	err := ps.Clean(rc, as)
	if err != nil {
		ps.Logger.Warnf("error while cleaning request, somehow: %+v", err)
	}
	status := fasthttp.StatusOK
	cutil.WriteCORS(rc)
	startNanos := time.Now().UnixNano()
	var redir string
	logger := ps.Logger
	ctx := ps.Context
	if !telemetry.SkipControllerMetrics {
		var span *telemetry.Span
		ctx, span, logger = telemetry.StartSpan(ps.Context, "controller."+key, ps.Logger)
		defer span.Complete()
	}
	logger = logger.With("path", string(rc.URI().Path()), "method", ps.Method, "status", status)
	ps.Context = ctx

	if ps.ForceRedirect == "" || ps.ForceRedirect == string(rc.URI().Path()) {
		redir, err = safeRun(f, as, ps)
		if err != nil {
			redir, err = handleError(key, as, ps, rc, err)
			if err != nil {
				ps.Logger.Warnf("unable to handle error: %+v", err)
			}
		}
	} else {
		redir = ps.ForceRedirect
	}
	if redir != "" {
		rc.Response.Header.Set("Location", redir)
		status = fasthttp.StatusFound
		rc.SetStatusCode(status)
	}
	elapsedMillis := float64((time.Now().UnixNano()-startNanos)/int64(time.Microsecond)) / float64(1000)
	defer ps.Close()
	rc.Response.Header.Set("Server-Timing", fmt.Sprintf("server:dur=%.3f", elapsedMillis))
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
