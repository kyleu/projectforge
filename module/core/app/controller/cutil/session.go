package cutil

import (
	"context"{{{ if .HasModule "user" }}}
	"time"{{{ end }}}

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/lib/user"{{{ if .HasModule "user" }}}
	usr "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

var initialIcons = {{{ if .HasModule "search" }}}[]string{"searchbox"}{{{ else }}}[]string{}{{{ end }}}

func LoadPageState(as *app.State, rc *fasthttp.RequestCtx, key string, logger util.Logger) *PageState {
	parentCtx, logger := httpmetrics.ExtractHeaders(rc, logger)
	ctx, span, logger := telemetry.StartSpan(parentCtx, "http:"+key, logger)
	span.Attribute("path", string(rc.Request.URI().Path()))
	if !telemetry.SkipControllerMetrics {
		httpmetrics.InjectHTTP(rc, span)
	}
	session, flashes, prof{{{ if .HasModule "oauth" }}}, accts{{{ end }}} := loadSession(ctx, as, rc, logger)
	params := ParamSetFromRequest(rc){{{ if .HasModule "oauth" }}}

	isAuthed, _ := user.Check("/", accts)
	isAdmin, _ := user.Check("/admin", accts){{{ end }}}{{{ if .HasModule "user" }}}

	u, _ := as.Services.User.Get(ctx, nil, prof.ID, logger){{{ end }}}

	return &PageState{
		Method: string(rc.Method()), URI: rc.Request.URI(), Flashes: flashes, Session: session,
		{{{ if .HasModule "user" }}}User: u, {{{ end }}}Profile: prof, {{{ if .HasModule "oauth" }}}Accounts: accts, Authed: isAuthed, Admin: isAdmin, {{{ end }}}Params: params,
		Icons: initialIcons, Context: ctx, Span: span, Logger: logger,
	}
}

func loadSession({{{ if .HasModule "user" }}}ctx{{{ else }}}_{{{ end }}} context.Context, {{{ if .HasModule "user" }}}as{{{ else }}}_{{{ end }}} *app.State, rc *fasthttp.RequestCtx, logger util.Logger) (util.ValueMap, []string, *user.Profile{{{ if .HasModule "oauth" }}}, user.Accounts{{{ end }}}) {
	sessionBytes := rc.Request.Header.Cookie(util.AppKey)
	session := util.ValueMap{}
	if len(sessionBytes) > 0 {
		dec, err := util.DecryptMessage(nil, string(sessionBytes), logger)
		if err != nil {
			logger.Warnf("error decrypting session: %+v", err)
		}
		err = util.FromJSON([]byte(dec), &session)
		if err != nil {
			session = util.ValueMap{}
		}
	}

	flashes := util.StringSplitAndTrim(session.GetStringOpt(csession.WebFlashKey), ";")
	if len(flashes) > 0 {
		delete(session, csession.WebFlashKey)
		err := csession.SaveSession(rc, session, logger)
		if err != nil {
			logger.Warnf("can't save session: %+v", err)
		}
	}

	prof, err := loadProfile(session)
	if err != nil {
		logger.Warnf("can't load profile: %+v", err)
	}{{{ if .HasModule "oauth" }}}

	var accts user.Accounts
	authX, ok := session[csession.WebAuthKey]
	if ok {
		authS, ok := authX.(string)
		if ok {
			accts = user.AccountsFromString(authS)
		}
	}{{{ end }}}{{{ if .HasModule "user" }}}

	if prof.ID == util.UUIDDefault {
		prof.ID = util.UUID()
		u := &usr.User{ID: prof.ID, Name: prof.Name{{{ if .HasModule "oauth" }}}, Picture: accts.Image(){{{ end }}}, Created: time.Now()}
		err = as.Services.User.Save(ctx, nil, logger, u)
		if err != nil {
			logger.Warnf("unable to save user [%s]", prof.ID.String())
			return nil, nil, prof, nil
		}
		session["profile"] = prof
		err = csession.SaveSession(rc, session, logger)
		if err != nil {
			logger.Warnf("unable to save session for user [%s]", prof.ID.String())
			return nil, nil, prof, nil
		}
	}{{{ end }}}

	return session, flashes, prof{{{ if .HasModule "oauth" }}}, accts{{{ end }}}
}

func loadProfile(session util.ValueMap) (*user.Profile, error) {
	x, ok := session["profile"]
	if !ok {
		return user.DefaultProfile.Clone(), nil
	}
	s, ok := x.(string)
	if !ok {
		m, ok := x.(map[string]any)
		if !ok {
			return user.DefaultProfile.Clone(), nil
		}
		s = util.ToJSON(m)
	}
	p := &user.Profile{}
	err := util.FromJSON([]byte(s), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
