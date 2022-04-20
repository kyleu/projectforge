package controller

import (
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/user"{{{ if.HasModule "marketing" }}}
	"{{{ .Package }}}/app/site"{{{ end }}}
)

const ({{{ if .HasModule "search" }}}
	defaultSearchPath  = "/search"{{{ end }}}
	defaultProfilePath = "/profile"
	defaultIcon        = "app"
)

func act(key string, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentAppState
	ps := loadPageState(rc, key, as)
	if allowed, reason := user.Check(string(ps.URI.Path()), ps.Accounts); !allowed {
		f = Unauthorized(rc, reason, ps.Accounts)
	}
	if err := initAppRequest(as, ps); err != nil {
		as.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, rc, f)
}
{{{ if.HasModule "marketing" }}}
func actSite(key string, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentSiteState
	ps := loadPageState(rc, key, as)
	ps.Menu = site.Menu(ps.Context, as, ps.Profile, ps.Accounts)
	if allowed, reason := user.Check(string(ps.URI.Path()), ps.Accounts); !allowed {
		f = Unauthorized(rc, reason, ps.Accounts)
	}
	if err := initSiteRequest(as, ps); err != nil {
		as.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, rc, f)
}
{{{ end }}}
func actComplete(key string, as *app.State, ps *cutil.PageState, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	err := clean(as, ps)
	if err != nil {
		as.Logger.Warnf("error while cleaning request, somehow: %+v", err)
	}
	status := fasthttp.StatusOK
	cutil.WriteCORS(rc)
	startNanos := time.Now().UnixNano()
	var redir string
	if ps.ForceRedirect == "" || ps.ForceRedirect == string(rc.URI().Path()) {
		redir, err = safeRun(f, as, ps)
		if err != nil {
			redir, err = handleError(key, as, ps, rc, err)
			if err != nil {
				as.Logger.Warnf("unable to handle error: %+v", err)
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
	l := ps.Logger.With(zap.String("path", string(rc.URI().Path())), zap.String("method", ps.Method), zap.Int("status", status), zap.Float64("elapsed", elapsedMillis))
	l.Debugf("processed request in [%.3fms] (render: %.3fms)", elapsedMillis, ps.RenderElapsed)
}

func safeRun(f func(as *app.State, ps *cutil.PageState) (string, error), as *app.State, ps *cutil.PageState) (s string, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			if recoverErr, ok := rec.(error); ok {
				e = errors.Wrap(recoverErr, "panic")
			} else {
				e = errors.Errorf("controller encountered panic recovery of type [%T]", rec)
			}
		}
	}()
	s, e = f(as, ps)
	return
}
