package cutil

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/mileusna/useragent"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/telemetry/httpmetrics"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

var (
	initialIcons = []string{"searchbox"}
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
	session, flashes, prof := loadSession(ctx, as, w, r, logger) //nolint:contextcheck
	params := ParamSetFromRequest(r)
	ua := useragent.Parse(r.Header.Get("User-Agent"))
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
	b, _ := io.ReadAll(http.MaxBytesReader(w, r.Body, MaxBodySize))
	r.Body = io.NopCloser(bytes.NewBuffer(b))

	return &PageState{
		Action: key, Method: r.Method, URI: r.URL, Flashes: flashes, Session: session,
		OS: os, OSVersion: ua.OSVersion, Browser: browser, BrowserVersion: ua.Version, Platform: platform,
		Profile: prof, Params: params,
		Icons: slices.Clone(initialIcons), Started: util.TimeCurrent(), Logger: logger, Context: ctx, Span: span, RequestBody: b, W: w,
	}
}

func loadSession(
	_ context.Context, _ *app.State, w http.ResponseWriter, r *http.Request, logger util.Logger,
) (util.ValueMap, []string, *user.Profile) {
	c, _ := r.Cookie(util.AppKey)
	if c == nil || c.Value == "" {
		return util.ValueMap{}, nil, user.DefaultProfile.Clone()
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
	p, err := util.FromJSONObj[*user.Profile]([]byte(s))
	if err != nil {
		return nil, err
	}
	if p.Name == "" {
		p.Name = user.DefaultProfile.Name
	}
	return p, nil
}
