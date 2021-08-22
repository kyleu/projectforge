package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/app/web"
	"{{{ .Package }}}/views"
	"{{{ .Package }}}/views/layout"
	"{{{ .Package }}}/views/verror"
)

var (
	_currentAppState  *app.State
	_currentSiteState *app.State
	initialIcons      = []string{"search"}
)

func SetAppState(a *app.State) {
	_currentAppState = a
	initApp(a)
}

func SetSiteState(a *app.State) {
	_currentSiteState = a
	initSite(a)
}

func rcRequiredString(rc *fasthttp.RequestCtx, key string, allowEmpty bool) (string, error) {
	v, ok := rc.UserValue(key).(string)
	if !ok || ((!allowEmpty) && v == "") {
		return v, errors.Errorf("must provide [%s] in path", key)
	}
	return v, nil
}

func render(rc *fasthttp.RequestCtx, as *app.State, page layout.Page, ps *cutil.PageState, breadcrumbs ...string) (string, error) {
	defer func() {
		x := recover()
		if x != nil {
			ps.Logger.Error(fmt.Sprintf("error processing template: %+v", x))
			switch t := x.(type) {
			case error:
				ed := util.GetErrorDetail(t)
				verror.WriteDetail(rc, ed, as, ps)
			default:
				ed := &util.ErrorDetail{Message: fmt.Sprint(t)}
				verror.WriteDetail(rc, ed, as, ps)
			}
		}
	}()
	ps.Breadcrumbs = append(ps.Breadcrumbs, breadcrumbs...)
	ct := cutil.GetContentType(rc)
	if ps.Data != nil {
		if cutil.IsContentTypeJSON(ct) {
			return cutil.RespondJSON(rc, "", ps.Data)
		} else if cutil.IsContentTypeXML(ct) {
			return cutil.RespondXML(rc, "", ps.Data)
		}
	}
	startNanos := time.Now().UnixNano()
	rc.Response.Header.SetContentType("text/html; charset=UTF-8")
	views.WriteRender(rc, page, as, ps)
	ps.RenderElapsed = float64((time.Now().UnixNano()-startNanos)/int64(time.Microsecond)) / float64(1000)
	return "", nil
}

func ersp(msg string, args ...interface{}) (string, error) {
	return "", errors.Errorf(msg, args...)
}

func flashAndRedir(success bool, msg string, redir string, rc *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	status := "error"
	if success {
		status = "success"
	}
	ps.Session.AddFlash(fmt.Sprintf("%s:%s", status, msg))
	if err := web.SaveSession(rc, ps.Session, ps.Logger); err != nil {
		return "", errors.Wrap(err, "unable to save flash session")
	}

	if strings.HasPrefix(redir, "/") {
		return redir, nil
	}
	if strings.HasPrefix(redir, "http") {
		ps.Logger.Warn("flash redirect attempted for non-local request")
		return "/", nil
	}
	return redir, nil
}
