package auth

import (
	"github.com/go-gem/sessions"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/user"
)

func getAuthURL(prv *Provider, rc *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) (string, error) {
	g, err := gothFor(rc, prv)
	if err != nil {
		return "", err
	}

	sess, err := g.BeginAuth(setState(rc))
	if err != nil {
		return "", err
	}

	u, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = cutil.StoreInSession(prv.ID, sess.Marshal(), rc, websess, logger)
	if err != nil {
		return "", err
	}

	return u, err
}

func getCurrentAuths(websess *sessions.Session) user.Accounts {
	authS, err := cutil.GetFromSession(WebSessKey, websess)
	var ret user.Accounts
	if err == nil && authS != "" {
		ret = user.AccountsFromString(authS)
	}
	return ret
}

func setCurrentAuths(s user.Accounts, rc *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	s.Sort()
	if len(s) == 0 {
		return cutil.RemoveFromSession(WebSessKey, rc, websess, logger)
	}
	return cutil.StoreInSession(WebSessKey, s.String(), rc, websess, logger)
}
