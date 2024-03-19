// Package controller - Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/verror"
)

func Render(w http.ResponseWriter, r *http.Request, as *app.State, page layout.Page, ps *cutil.PageState, breadcrumbs ...string) (string, error) {
	defer func() {
		x := recover()
		if x != nil {
			ps.LogError("error processing template: %+v", x)
			switch t := x.(type) {
			case error:
				ed := util.GetErrorDetail(t, ps.Admin)
				verror.WriteDetail(w, ed, as, ps)
			default:
				ed := &util.ErrorDetail{Type: fmt.Sprintf("%T", x), Message: fmt.Sprint(t)}
				verror.WriteDetail(w, ed, as, ps)
			}
		}
	}()
	ps.Breadcrumbs = append(ps.Breadcrumbs, breadcrumbs...)
	ct := cutil.GetContentType(r)
	if ps.Data != nil {
		var fn string
		if r.URL.Query().Get("download") == "true" {
			fn = ps.Action
		}
		switch {
		case cutil.IsContentTypeCSV(ct):
			return cutil.RespondCSV(w, fn, ps.Data)
		case cutil.IsContentTypeJSON(ct):
			return cutil.RespondJSON(w, fn, ps.Data)
		case cutil.IsContentTypeXML(ct):
			return cutil.RespondXML(w, fn, ps.Data)
		case cutil.IsContentTypeYAML(ct):
			return cutil.RespondYAML(w, fn, ps.Data)
		case cutil.IsContentTypeDebug(ct):
			return cutil.RespondDebug(w, r, as, fn, ps)
		}
	}
	startNanos := util.TimeCurrentNanos()
	w.Header().Set(cutil.HeaderContentType, "text/html; charset=UTF-8")
	views.WriteRender(w, page, as, ps)
	ps.RenderElapsed = float64((util.TimeCurrentNanos()-startNanos)/int64(time.Microsecond)) / float64(1000)
	return "", nil
}

func ERsp(msg string, args ...any) (string, error) {
	return "", errors.Errorf(msg, args...)
}

func FlashAndRedir(success bool, msg string, redir string, w http.ResponseWriter, ps *cutil.PageState) (string, error) {
	status := util.KeyError
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
	if err := csession.SaveSession(w, ps.Session, ps.Logger); err != nil {
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

func ReturnToReferrer(msg string, dflt string, w http.ResponseWriter, ps *cutil.PageState) (string, error) {
	refer := ""
	referX, ok := ps.Session[csession.ReferKey]
	if ok {
		refer, ok = referX.(string)
		if ok {
			_ = csession.RemoveFromSession(csession.ReferKey, w, ps.Session, ps.Logger)
		}
	}
	if refer == "" {
		refer = dflt
	}
	return FlashAndRedir(true, msg, refer, w, ps)
}
