package cutil

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/mileusna/useragent"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/lib/user"{{{ if .DatabaseUISaveUser }}}
	usr "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

var (
	initialIcons = {{{ if .HasModule "search" }}}[]string{"searchbox"}{{{ else }}}[]string{}{{{ end }}}
	MaxBodySize  = int64(1024 * 1024 * 128) // 128MB
)

// @contextcheck(req_has_ctx).
func LoadPageState(as *app.State, w *WriteCounter, r *http.Request, key string, logger util.Logger) *PageState {
	parentCtx, logger := httpmetrics.ExtractHeaders(r, logger)
	ctx, span, logger := telemetry.StartSpan(parentCtx, "http:"+key, logger)
	span.Attribute("path", r.URL.Path)
	if !telemetry.SkipControllerMetrics {
		httpmetrics.InjectHTTP(200, r, span)
	}
	session, flashes, prof{{{ if .HasAccount }}}, accts{{{ end }}} := loadSession(ctx, as, w, r, logger) //nolint:contextcheck
	params := ParamSetFromRequest(r)
	ua := useragent.Parse(r.Header.Get("User-Agent"))
	os := strings.ToLower(ua.OS)
	browser := strings.ToLower(ua.Name)
	platform := util.KeyUnknown
	switch {
	case ua.Desktop:
		platform = "desktop"
	case ua.Tablet:
		platform = "tablet"
	case ua.Mobile:
		platform = "mobile"
	case ua.Bot:
		platform = "bot"
	}
	span.Attribute("browser", browser)
	span.Attribute("os", os){{{ if .HasAccount }}}

	isAuthed, _ := user.Check("/", accts)
	isAdmin, _ := user.Check("/admin", accts){{{ end }}}
	b, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, MaxBodySize))
	r.Body = io.NopCloser(bytes.NewBuffer(b)){{{ if .HasUser }}}

	u, _ := as.User(ctx, prof.ID, logger){{{ end }}}

	return &PageState{
		Action: key, Method: r.Method, URI: r.URL, Flashes: flashes, Session: session,
		OS: os, OSVersion: ua.OSVersion, Browser: browser, BrowserVersion: ua.Version, Platform: platform, Transport: r.URL.Scheme,
		{{{ if .HasUser }}}User: u, {{{ end }}}Profile: prof, {{{ if .HasAccount }}}Accounts: accts, Authed: isAuthed, Admin: isAdmin, {{{ end }}}Params: params,
		Icons: util.ArrayCopy(initialIcons), Started: util.TimeCurrent(), Logger: logger, Context: ctx, Span: span, RequestBody: b, W: w,
	}
}

func loadSession(
	{{{ if .DatabaseUISaveUser }}}ctx{{{ else }}}_{{{ end }}} context.Context, {{{ if .DatabaseUISaveUser }}}as{{{ else }}}_{{{ end }}} *app.State, w http.ResponseWriter, r *http.Request, logger util.Logger,
) (util.ValueMap, []string, *user.Profile{{{ if .HasAccount }}}, user.Accounts{{{ end }}}) {
	c, _ := r.Cookie(util.AppKey)
	if c == nil || c.Value == "" {
		return util.ValueMap{}, nil, user.DefaultProfile.Clone(){{{ if .HasAccount }}}, nil{{{ end }}}
	}

	dec, err := util.DecryptMessage(nil, c.Value, logger)
	if err != nil {
		logger.Warnf("error decrypting session: %+v", err)
	}
	session, err := util.FromJSONMap([]byte(dec))
	if err != nil {
		session = util.ValueMap{}
	}

	flashes := util.StringSplitAndTrim(session.GetStringOpt(csession.WebFlashKey), ";")
	if len(flashes) > 0 {
		delete(session, csession.WebFlashKey)
		if e := csession.SaveSession(w, session, logger); e != nil {
			logger.Warnf("can't save session: %+v", e)
		}
	}

	prof, err := loadProfile(session)
	if err != nil {
		logger.Warnf("can't load profile: %+v", err)
	}{{{ if .HasAccount }}}

	var accts user.Accounts
	authX, ok := session[csession.WebAuthKey]
	if ok {
		authS, ok := authX.(string)
		if ok {
			accts = user.AccountsFromString(authS)
		}
	}{{{ end }}}{{{ if .HasUser }}}

	if prof.ID == util.UUIDDefault {
		prof.ID = util.UUID(){{{ if .DatabaseUISaveUser }}}
		u := &usr.User{ID: prof.ID, Name: prof.Name{{{ if .HasAccount }}}, Picture: accts.Image(){{{ end }}}, Created: util.TimeCurrent()}
		err = as.Services.User.Save(ctx, nil, logger, u)
		if err != nil {
			logger.Warnf("unable to save user [%s]", prof.ID.String())
			return nil, nil, prof, nil
		}{{{ end }}}
		session["profile"] = prof
		err = csession.SaveSession(w, session, logger)
		if err != nil {
			logger.Warnf("unable to save session for user [%s]", prof.ID.String())
			return nil, nil, prof{{{ if .HasAccount }}}, nil{{{ end }}}
		}
	}{{{ end }}}

	return session, flashes, prof{{{ if .HasAccount }}}, accts{{{ end }}}
}

func loadProfile(session util.ValueMap) (*user.Profile, error) {
	x, ok := session["profile"]
	if !ok {
		return user.DefaultProfile.Clone(), nil
	}
	s, err := util.Cast[string](x)
	if err != nil {
		m, err := util.Cast[map[string]any](x)
		if err != nil {
			return user.DefaultProfile.Clone(), nil //nolint:nilerr
		}
		s = util.ToJSON(m)
	}
	p, err := util.FromJSONObj[*user.Profile]([]byte(s))
	if err != nil {
		return nil, err
	}
	if p.Name == "" {
		p.Name = user.DefaultProfile.Name
	}
	return p, nil
}
