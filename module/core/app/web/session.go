package web

import (
	"github.com/go-gem/sessions"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	WebSessKey = "auth"
	ReferKey   = "refer"
)

var webSessOpts = &sessions.Options{Path: "/", HttpOnly: true, MaxAge: 365 * 24 * 60 * 60 /* , SameSite: http.SameSiteStrictMode */}

func StoreInSession(k string, v string, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	websess.Values[k] = v
	return SaveSession(ctx, websess, logger)
}

func RemoveFromSession(k string, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	delete(websess.Values, k)
	return SaveSession(ctx, websess, logger)
}

func SaveSession(ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	// logger.Debugf("saving session with [%d] keys", len(websess.Values))
	websess.Options = webSessOpts
	return websess.Save(ctx)
}

func GetFromSession(key string, websess *sessions.Session) (string, error) {
	value, ok := websess.Values[key]
	if !ok {
		return "", errors.Errorf("could not find a matching session value with key [%s] for this request", key)
	}
	s, ok := value.(string)
	if !ok {
		return "", errors.Errorf("session value with key [%s] is of type [%T], not [string]", key, value)
	}
	return s, nil
}
