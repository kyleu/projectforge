package controller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/kyleu/projectforge/app/auth"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views"
	"github.com/kyleu/projectforge/views/verror"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/go-gem/sessions"
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/views/layout"
)

var (
	_currentAppState  *app.State
	_currentSiteState *app.State
	_rootLogger       *zap.SugaredLogger
	initialIcons      = []string{"search"}
)

var sessionKey = func() string {
	x := os.Getenv("SESSION_KEY")
	if x == "" {
		x = util.AppKey + "_random_secret_key"
	}
	return x
}()

var store *sessions.CookieStore

func SetAppState(a *app.State, l *zap.SugaredLogger) {
	_currentAppState = a
	_rootLogger = l
	initApp()
}

func SetSiteState(a *app.State, l *zap.SugaredLogger) {
	_currentSiteState = a
	_rootLogger = l
	initSite()
}

func ctxRequiredString(ctx *fasthttp.RequestCtx, key string, allowEmpty bool) (string, error) {
	v, ok := ctx.UserValue(key).(string)
	if !ok || ((!allowEmpty) && v == "") {
		return v, errors.Errorf("must provide [%s] in path", key)
	}
	return v, nil
}

func render(ctx *fasthttp.RequestCtx, as *app.State, page layout.Page, ps *cutil.PageState, breadcrumbs ...string) (string, error) {
	defer func() {
		x := recover()
		if x != nil {
			ps.Logger.Error(fmt.Sprintf("error processing template: %+v", x))
			switch t := x.(type) {
			case error:
				ed := util.GetErrorDetail(t)
				verror.WriteDetail(ctx, ed, as, ps)
			default:
				ed := &util.ErrorDetail{Message: fmt.Sprintf("%v", t)}
				verror.WriteDetail(ctx, ed, as, ps)
			}
		}
	}()
	ps.Breadcrumbs = append(ps.Breadcrumbs, breadcrumbs...)
	ct := cutil.GetContentType(ctx)
	if ps.Data != nil {
		if cutil.IsContentTypeJSON(ct) {
			return cutil.RespondJSON(ctx, "", ps.Data)
		} else if cutil.IsContentTypeXML(ct) {
			return cutil.RespondXML(ctx, "", ps.Data)
		}
	}
	startNanos := time.Now().UnixNano()
	ctx.Response.Header.SetContentType("text/html; charset=UTF-8")
	views.WriteRender(ctx, page, as, ps)
	ps.RenderElapsed = float64((time.Now().UnixNano()-startNanos)/int64(time.Microsecond)) / float64(1000)
	return "", nil
}

func ersp(msg string, args ...interface{}) (string, error) {
	return "", errors.Errorf(msg, args...)
}

func flashAndRedir(success bool, msg string, redir string, ctx *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	status := "error"
	if success {
		status = "success"
	}
	ps.Session.AddFlash(fmt.Sprintf("%s:%s", status, msg))
	if err := auth.SaveSession(ctx, ps.Session, ps.Logger); err != nil {
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

func initStore(keyPairs ...[]byte) *sessions.CookieStore {
	ret := sessions.NewCookieStore(keyPairs...)
	for _, x := range ret.Codecs {
		c, ok := x.(*securecookie.SecureCookie)
		if ok {
			c.MaxLength(65536)
		}
	}
	return ret
}
