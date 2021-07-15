package auth

import (
	"github.com/go-gem/sessions"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const ReferKey = "auth-refer"

func BeginAuthHandler(prv *Provider, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) (string, error) {
	u, err := getAuthURL(prv, ctx, websess, logger)
	if err != nil {
		return "", err
	}
	refer := string(ctx.Request.URI().QueryArgs().Peek("refer"))
	if refer != "" && refer != "/profile" {
		_ = StoreInSession(ReferKey, refer, ctx, websess, logger)
	}
	return u, nil
}

func CompleteUserAuth(prv *Provider, ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) (*Session, Sessions, error) {
	value, err := getFromSession(prv.ID, websess)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = removeProviderData(ctx, websess, logger)
	}()

	g, err := gothFor(ctx, prv)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create oauth provider")
	}

	sess, err := g.UnmarshalSession(value)
	if err != nil {
		return nil, nil, err
	}

	err = validateState(ctx, sess)
	if err != nil {
		return nil, nil, err
	}

	user, err := g.FetchUser(sess)
	if err == nil {
		return addToSession(user.Provider, user.Email, ctx, websess, logger)
	}

	_, err = sess.Authorize(g, &params{q: ctx.Request.URI().QueryArgs()})
	if err != nil {
		return nil, nil, err
	}

	err = StoreInSession(prv.ID, sess.Marshal(), ctx, websess, logger)
	if err != nil {
		return nil, nil, err
	}

	gu, err := g.FetchUser(sess)
	if err != nil {
		return nil, nil, err
	}

	return addToSession(gu.Provider, gu.Email, ctx, websess, logger)
}

func gothFor(ctx *fasthttp.RequestCtx, prv *Provider) (goth.Provider, error) {
	proto := string(ctx.URI().Scheme())
	host := string(ctx.URI().Host())
	if host == "" {
		host = "localhost"
	}
	return prv.Goth(proto, host)
}

func Logout(ctx *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger, prvKeys ...string) error {
	a := getCurrentAuths(websess)
	n := a.Purge(prvKeys...)
	return setCurrentAuths(n, ctx, websess, logger)
}
