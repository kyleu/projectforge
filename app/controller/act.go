package controller

import (
	"fmt"
	"time"

	"github.com/kyleu/projectforge/app/auth"
	"github.com/kyleu/projectforge/app/site"
	"github.com/kyleu/projectforge/app/theme"
	"github.com/valyala/fasthttp"

	"go.uber.org/zap"

	"github.com/kyleu/projectforge/views/verror"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/util"
)

const (
	defaultSearchPath  = "/search"
	defaultProfilePath = "/profile"
	defaultKey         = "app"
	defaultIcon        = defaultKey
)

func act(key string, ctx *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	mode := defaultKey
	as, ps := actPrepare(mode, ctx)
	clean(as, ps)
	actComplete(key, as, ps, ctx, f)
}

func actSite(key string, ctx *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	as, ps := actPrepare("site", ctx)
	ps.Menu = site.Menu(ps.Profile, ps.Auth)
	clean(as, ps)
	actComplete(key, as, ps, ctx, f)
}

func actPrepare(mode string, ctx *fasthttp.RequestCtx) (*app.State, *cutil.PageState) {
	if mode == "site" {
		return _currentSiteState, loadPageState(ctx)
	}
	return _currentAppState, loadPageState(ctx)
}

func loadPageState(ctx *fasthttp.RequestCtx) *cutil.PageState {
	logger := _rootLogger.With(zap.String("path", string(ctx.Request.URI().Path())))

	if store == nil {
		store = initStore([]byte(sessionKey))
	}
	session, err := store.Get(ctx, util.AppKey)
	if err != nil {
		logger.Warnf("error retrieving session: %+v", err)
		session, err = store.New(ctx, util.AppKey)
		if err != nil {
			logger.Warnf("error creating new session: %+v", err)
		}
	}
	flashes := util.StringArrayFromInterfaces(session.Flashes())
	if len(flashes) > 0 {
		err = auth.SaveSession(ctx, session, logger)
		if err != nil {
			logger.Warnf("can't save session: %+v", err)
		}
	}

	prof, err := loadProfile(session)
	if err != nil {
		logger.Warnf("can't load profile: %+v", err)
	}

	var a auth.Sessions
	authX, ok := session.Values["auth"]
	if ok {
		authS, ok := authX.(string)
		if ok {
			a = auth.SessionsFromString(authS)
		}
	}

	return &cutil.PageState{
		Method:  string(ctx.Method()),
		URI:     ctx.Request.URI(),
		Flashes: flashes,
		Session: session,
		Profile: prof,
		Auth:    a,
		Icons:   initialIcons,
		Logger:  logger,
	}
}

func actComplete(key string, as *app.State, ps *cutil.PageState, ctx *fasthttp.RequestCtx, f func(as *app.State, ps *cutil.PageState) (string, error)) {
	status := fasthttp.StatusOK
	cutil.WriteCORS(ctx)
	startNanos := time.Now().UnixNano()
	redir, err := f(as, ps)
	if err != nil {
		status = fasthttp.StatusInternalServerError
		ctx.SetStatusCode(status)

		ps.Logger.Errorf("error running action [%s]: %+v", key, err)

		if len(ps.Breadcrumbs) == 0 {
			ps.Breadcrumbs = cutil.Breadcrumbs{"Error"}
		}
		errDetail := util.GetErrorDetail(err)
		page := &verror.Error{Err: errDetail}

		clean(as, ps)
		redir, err = render(ctx, as, page, ps)
		if err != nil {
			_, _ = ctx.Write([]byte(fmt.Sprintf("error while running error handler: %+v", err)))
		}
	}
	if redir != "" {
		ctx.Response.Header.Set("Location", redir)
		status = fasthttp.StatusFound
		ctx.SetStatusCode(status)
	}
	elapsedMillis := float64((time.Now().UnixNano()-startNanos)/int64(time.Microsecond)) / float64(1000)
	l := ps.Logger.With(zap.String("method", ps.Method), zap.Int("status", status), zap.Float64("elapsed", elapsedMillis))
	l.Debugf("processed request in [%.3fms] (render: %.3fms)", elapsedMillis, ps.RenderElapsed)
}

func clean(as *app.State, ps *cutil.PageState) {
	if ps.Profile != nil && ps.Profile.Theme == "" {
		ps.Profile.Theme = theme.ThemeDefault.Key
	}
	if ps.RootIcon == "" {
		ps.RootIcon = defaultIcon
	}
	if ps.RootPath == "" {
		ps.RootPath = "/"
	}
	if ps.RootTitle == "" {
		ps.RootTitle = util.AppName
	}
	if ps.SearchPath == "" {
		ps.SearchPath = defaultSearchPath
	}
	if ps.ProfilePath == "" {
		ps.ProfilePath = defaultProfilePath
	}
	if len(ps.Menu) == 0 {
		ps.Menu = MenuFor(as)
	}
}
