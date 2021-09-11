package controller

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package}}}/app"
	"{{{ .Package}}}/app/controller/cutil"
	"{{{ .Package}}}/app/telemetry"
	"{{{ .Package}}}/app/telemetry/httpmetrics"
	"{{{ .Package}}}/app/user"
	"{{{ .Package}}}/app/util"
)

func loadPageState(rc *fasthttp.RequestCtx, key string, as *app.State) *cutil.PageState {
	rc = httpmetrics.ExtractHeaders(rc, as.Logger)
	traceCtx, span := telemetry.StartSpan(rc, "pagestate", "http:"+key)
	httpmetrics.InjectHTTP(rc, span)

	path := string(rc.Request.URI().Path())
	sc := span.SpanContext()
	logger := as.Logger.With(zap.String("path", path), zap.String("trace", sc.TraceID().String()), zap.String("span", sc.SpanID().String()))

	session, flashes, prof, accts := loadSession(rc, logger)

	isAuthed, _ := user.Check("/", accts)
	isAdmin, _ := user.Check("/admin", accts)

	return &cutil.PageState{
		Method:   string(rc.Method()),
		URI:      rc.Request.URI(),
		Flashes:  flashes,
		Session:  session,
		Profile:  prof,
		Accounts: accts,
		Authed:   isAuthed,
		Admin:    isAdmin,
		Icons:    initialIcons,
		Context:  traceCtx,
		Span:     &span,
		Logger:   logger,
	}
}

func loadSession(rc *fasthttp.RequestCtx, logger *zap.SugaredLogger) (util.ValueMap, []string, *user.Profile, user.Accounts) {
	sessionBytes := rc.Request.Header.Cookie(util.AppKey)
	session := util.ValueMap{}
	if len(sessionBytes) > 0 {
		dec, err := cutil.DecryptMessage(string(sessionBytes), logger)
		if err != nil {
			logger.Warnf("error decrypting session: %+v", err)
		}
		err = util.FromJSON([]byte(dec), &session)
		if err != nil {
			logger.Warnf("error parsing session: %+v", err)
		}
	}

	flashes := util.SplitAndTrim(session.GetStringOpt(cutil.WebFlashKey), ",")
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
