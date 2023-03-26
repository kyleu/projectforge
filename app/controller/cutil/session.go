// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"context"

	"github.com/valyala/fasthttp"

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

	return &PageState{
		Method: string(rc.Method()), URI: rc.Request.URI(), Flashes: flashes, Session: session,
		Profile: prof, Params: params,
		Icons: initialIcons, Context: ctx, Span: span, Logger: logger,
	}
}

func loadSession(ctx context.Context, as *app.State, rc *fasthttp.RequestCtx, logger util.Logger) (util.ValueMap, []string, *user.Profile) {
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
