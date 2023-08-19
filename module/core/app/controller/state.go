package controller

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/verror"
)

var (
	{{{ if.HasModule "marketing" }}}_currentAppState       *app.State
	_currentAppRootLogger  util.Logger
	_currentSiteState      *app.State
	_currentSiteRootLogger util.Logger{{{ else }}}_currentAppState      *app.State
	_currentAppRootLogger util.Logger{{{ end }}}
)

func SetAppState(a *app.State, logger util.Logger) {
	_currentAppState = a
	_currentAppRootLogger = logger
	initApp(a, logger)
}{{{ if.HasModule "marketing" }}}

func SetSiteState(a *app.State, logger util.Logger) {
	_currentSiteState = a
	_currentSiteRootLogger = logger
	initSite(a, logger)
}{{{ end }}}

func handleError(key string, as *app.State, ps *cutil.PageState, rc *fasthttp.RequestCtx, err error) (string, error) {
	rc.SetStatusCode(fasthttp.StatusInternalServerError)

	ps.LogError("error running action [%s]: %+v", key, err)

	if len(ps.Breadcrumbs) == 0 {
		bc := util.StringSplitAndTrim(string(rc.URI().Path()), "/")
		bc = append(bc, "Error")
		ps.Breadcrumbs = bc
	}

	if cleanErr := ps.Clean(rc, as); cleanErr != nil {
		ps.Logger.Error(cleanErr)
		msg := fmt.Sprintf("error while cleaning request: %+v", cleanErr)
		ps.Logger.Error(msg)
		_, _ = rc.WriteString(msg)
		return "", cleanErr
	}

	e := util.GetErrorDetail(err)
	ps.Data = e
	redir, renderErr := Render(rc, as, &verror.Error{Err: e}, ps)
	if renderErr != nil {
		msg := fmt.Sprintf("error while running error handler: %+v", renderErr)
		ps.Logger.Error(msg)
		_, _ = rc.WriteString(msg)
		return "", renderErr
	}
	return redir, nil
}
