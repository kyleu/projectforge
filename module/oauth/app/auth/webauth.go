package auth

import (
	"github.com/go-gem/sessions"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/web"
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

	err = web.StoreInSession(prv.ID, sess.Marshal(), ctx, websess, logger)
	if err != nil {
		return "", err
	}

	return u, err
}

func getCurrentAuths(websess *sessions.Session) web.Accounts {
	authS, err := web.GetFromSession(WebSessKey, websess)
	var ret web.Accounts
	if err == nil && authS != "" {
		ret = web.AccountsFromString(authS)
	}
	return ret
}

func setCurrentAuths(s web.Accounts, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	s.Sort()
	if len(s) == 0 {
		return web.RemoveFromSession(WebSessKey, ctx, websess, logger)
	}
	return web.StoreInSession(WebSessKey, s.String(), ctx, websess, logger)
}