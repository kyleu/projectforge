package auth

import (
	"github.com/go-gem/sessions"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func getAuthURL(prv *Provider, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) (string, error) {
	g, err := gothFor(ctx, prv)
	if err != nil {
		return "", err
	}

	sess, err := g.BeginAuth(setState(ctx))
	if err != nil {
		return "", err
	}

	u, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = StoreInSession(prv.ID, sess.Marshal(), ctx, websess, logger)
	if err != nil {
		return "", err
	}

	return u, err
}

func getCurrentAuths(websess *sessions.Session) Sessions {
	authS, err := getFromSession(WebSessKey, websess)
	var ret Sessions
	if err == nil && authS != "" {
		ret = SessionsFromString(authS)
	}
	return ret
}

func setCurrentAuths(s Sessions, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	s.Sort()
	if len(s) == 0 {
		return RemoveFromSession(WebSessKey, ctx, websess, logger)
	}
	return StoreInSession(WebSessKey, s.String(), ctx, websess, logger)
}
