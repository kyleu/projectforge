package controller

import (
	"context"
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/verror"
)

var (
	_currentAppState       *app.State
	_currentAppRootLogger  util.Logger
	_currentSiteState      *app.State
	_currentSiteRootLogger util.Logger
)

func SetAppState(a *app.State, logger util.Logger) error {
	_currentAppState = a
	_currentAppRootLogger = logger
	return initApp(context.Background(), a, logger)
}

func SetSiteState(a *app.State, logger util.Logger) error {
	_currentSiteState = a
	_currentSiteRootLogger = logger
	return initSite(context.Background(), a, logger)
}

func handleError(key string, as *app.State, ps *cutil.PageState, r *http.Request, err error) (string, error) {
	ps.W.WriteHeader(http.StatusInternalServerError)

	ps.LogError("error running action [%s]: %+v", key, err)

	if len(ps.Breadcrumbs) == 0 {
		bc := util.StringSplitAndTrim(r.URL.Path, "/")
		bc = append(bc, "Error")
		ps.Breadcrumbs = bc
	}

	if cleanErr := ps.Clean(r, as); cleanErr != nil {
		ps.Logger.Error(cleanErr)
		msg := fmt.Sprintf("error while cleaning request: %+v", cleanErr)
		ps.Logger.Error(msg)
		_, _ = ps.W.Write([]byte(msg))
		return "", cleanErr
	}

	e := util.GetErrorDetail(err, ps.Admin)
	ps.Data = e
	if ps.Title == "" {
		ps.Title = "Error Encountered"
	}
	redir, renderErr := Render(r, as, &verror.Error{Err: e}, ps)
	if renderErr != nil {
		msg := fmt.Sprintf("error while running error handler: %+v", renderErr)
		ps.Logger.Error(msg)
		_, _ = ps.W.Write([]byte(msg))
		return "", renderErr
	}
	return redir, nil
}
