// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/verror"
)

func Render(rc *fasthttp.RequestCtx, as *app.State, page layout.Page, ps *cutil.PageState, breadcrumbs ...string) (string, error) {
	defer func() {
		x := recover()
		if x != nil {
			ps.Logger.Errorf("error processing template: %+v", x)
			switch t := x.(type) {
			case error:
				ed := util.GetErrorDetail(t)
				verror.WriteDetail(rc, ed, as, ps)
			default:
				ed := &util.ErrorDetail{Type: fmt.Sprintf("%T", x), Message: fmt.Sprint(t)}
				verror.WriteDetail(rc, ed, as, ps)
			}
		}
	}()
	ps.Breadcrumbs = append(ps.Breadcrumbs, breadcrumbs...)
	ct := cutil.GetContentType(rc)
	if ps.Data != nil {
		switch {
		case cutil.IsContentTypeJSON(ct):
			return cutil.RespondJSON(rc, "", ps.Data)
		case cutil.IsContentTypeXML(ct):
			return cutil.RespondXML(rc, "", ps.Data)
		case cutil.IsContentTypeYAML(ct):
			return cutil.RespondYAML(rc, "", ps.Data)
		case ct == "debug":
			return cutil.RespondDebug(rc, "", ps.Data)
		}
	}
	startNanos := time.Now().UnixNano()
	rc.Response.Header.SetContentType("text/html; charset=UTF-8")
	views.WriteRender(rc, page, as, ps)
	ps.RenderElapsed = float64((time.Now().UnixNano()-startNanos)/int64(time.Microsecond)) / float64(1000)
	return "", nil
}

func ERsp(msg string, args ...any) (string, error) {
	return "", errors.Errorf(msg, args...)
}

func FlashAndRedir(success bool, msg string, redir string, rc *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	status := "error"
	if success {
		status = "success"
	}
	msgFmt := fmt.Sprintf("%s:%s", status, msg)
	currStr := ps.Session.GetStringOpt(csession.WebFlashKey)
	if currStr == "" {
		currStr = msgFmt
	} else {
		curr := util.StringSplitAndTrim(currStr, ";")
		curr = append(curr, msgFmt)
		currStr = strings.Join(curr, ";")
	}
	ps.Session[csession.WebFlashKey] = currStr
	if err := csession.SaveSession(rc, ps.Session, ps.Logger); err != nil {
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

func ReturnToReferrer(msg string, dflt string, rc *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	refer := ""
	referX, ok := ps.Session[csession.ReferKey]
	if ok {
		refer, ok = referX.(string)
		if ok {
			_ = csession.RemoveFromSession(csession.ReferKey, rc, ps.Session, ps.Logger)
		}
	}
	if refer == "" {
		refer = dflt
	}
	return FlashAndRedir(true, msg, refer, rc, ps)
}
