package controller

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/lib/telemetry"
	"github.com/kyleu/projectforge/app/lib/telemetry/httpmetrics"
	"github.com/kyleu/projectforge/app/lib/user"
	"github.com/kyleu/projectforge/app/util"
)

func loadPageState(rc *fasthttp.RequestCtx, key string, as *app.State) *cutil.PageState {
	rc = httpmetrics.ExtractHeaders(rc, as.Logger)
	traceCtx, span := telemetry.StartSpan(rc, "pagestate", "http:"+key)
	httpmetrics.InjectHTTP(rc, span)

	path := string(rc.Request.URI().Path())
	sc := span.SpanContext()
	tid, sid := sc.TraceID().String(), sc.SpanID().String()
	logger := as.Logger.With(
		zap.String("path", path),
		zap.String("trace", tid), zap.String("dd.trace_id", tid),
		zap.String("span", sid), zap.String("dd.span_id", sid),
	)

	session, flashes, prof, accts := loadSession(rc, logger)

	isAuthed, _ := user.Check("/", accts)
	isAdmin, _ := user.Check("/admin", accts)

	return &cutil.PageState{
		Method: string(rc.Method()), URI: rc.Request.URI(), Flashes: flashes, Session: session,
		Profile: prof, Accounts: accts, Authed: isAuthed, Admin: isAdmin,
		Icons: initialIcons, Context: traceCtx, Span: span, Logger: logger,
	}
}

func loadSession(rc *fasthttp.RequestCtx, logger *zap.SugaredLogger) (util.ValueMap, []string, *user.Profile, user.Accounts) {
	sessionBytes := rc.Request.Header.Cookie(util.AppKey)
	session := util.ValueMap{}
	if len(sessionBytes) > 0 {
		dec, err := util.DecryptMessage(nil, string(sessionBytes), logger)
		if err != nil {
			logger.Warnf("error decrypting session: %+v", err)
		}
		err = util.FromJSON([]byte(dec), &session)
		if err != nil {
			logger.Warnf("error parsing session: %+v", err)
		}
	}

	flashes := util.StringSplitAndTrim(session.GetStringOpt(cutil.WebFlashKey), ",")
	if len(flashes) > 0 {
		delete(session, cutil.WebFlashKey)
		err := cutil.SaveSession(rc, session, logger)
		if err != nil {
			logger.Warnf("can't save session: %+v", err)
		}
	}

	prof, err := loadProfile(session)
	if err != nil {
		logger.Warnf("can't load profile: %+v", err)
	}

	var accts user.Accounts
	authX, ok := session[cutil.WebAuthKey]
	if ok {
		authS, ok := authX.(string)
		if ok {
			accts = user.AccountsFromString(authS)
		}
	}

	return session, flashes, prof, accts
}
