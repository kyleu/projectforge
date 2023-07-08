// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"context"
	"strings"

	"github.com/mileusna/useragent"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/telemetry/httpmetrics"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

var initialIcons = []string{"searchbox"}

func LoadPageState(as *app.State, rc *fasthttp.RequestCtx, key string, logger util.Logger) *PageState {
	parentCtx, logger := httpmetrics.ExtractHeaders(rc, logger)
	ctx, span, logger := telemetry.StartSpan(parentCtx, "http:"+key, logger)
	span.Attribute("path", string(rc.Request.URI().Path()))
	if !telemetry.SkipControllerMetrics {
		httpmetrics.InjectHTTP(rc, span)
	}
	session, flashes, prof := loadSession(ctx, as, rc, logger)
	params := ParamSetFromRequest(rc)
	ua := useragent.Parse(string(rc.Request.Header.Peek("User-Agent")))
	os := strings.ToLower(ua.OS)
	browser := strings.ToLower(ua.Name)
	platform := "unknown"
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
	span.Attribute("os", os)

	return &PageState{
		Method: string(rc.Method()), URI: rc.Request.URI(), Flashes: flashes, Session: session,
		OS: os, OSVersion: ua.OSVersion, Browser: browser, BrowserVersion: ua.Version, Platform: platform,
		Profile: prof, Params: params,
		Icons: slices.Clone(initialIcons), Started: util.TimeCurrent(), Logger: logger, Context: ctx, Span: span,
	}
}

func loadSession(_ context.Context, _ *app.State, rc *fasthttp.RequestCtx, logger util.Logger) (util.ValueMap, []string, *user.Profile) {
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
	}

	return session, flashes, prof
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
