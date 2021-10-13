package controller

import (
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/site"
	"github.com/kyleu/projectforge/app/user"
)

const (
	defaultSearchPath  = "/search"
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

func actSite(key string, rc *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as := _currentSiteState
	ps := loadPageState(rc, key, as)
	ps.Menu = site.Menu(ps.Profile, ps.Accounts)
	if allowed, reason := user.Check(string(ps.URI.Path()), ps.Accounts); !allowed {
		f = Unauthorized(rc, reason, ps.Accounts)
	}
	if err := initSiteRequest(as, ps); err != nil {
		as.Logger.Warnf("%+v", err)
	}
	actComplete(key, as, ps, rc, f)
}

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
		redir, err = f(as, ps)
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
	l := ps.Logger.With(zap.String("method", ps.Method), zap.Int("status", status), zap.Float64("elapsed", elapsedMillis))
	l.Debugf("processed request in [%.3fms] (render: %.3fms)", elapsedMillis, ps.RenderElapsed)
}
