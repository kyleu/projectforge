package auth

import (
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const defaultProfilePath = "/profile"

func BeginAuthHandler(prv *Provider, rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) (string, error) {
	u, err := getAuthURL(prv, rc, websess, logger)
	if err != nil {
		return "", err
	}
	refer := string(rc.Request.URI().QueryArgs().Peek("refer"))
	if refer != "" && refer != defaultProfilePath {
		_ = cutil.StoreInSession(cutil.ReferKey, refer, rc, websess, logger)
	}
	return u, nil
}

func CompleteUserAuth(prv *Provider, rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) (*user.Account, user.Accounts, error) {
	value, err := cutil.GetFromSession(prv.ID, websess)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = removeProviderData(rc, websess, logger)
	}()

	g, err := gothFor(rc, prv)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create oauth provider")
	}

	sess, err := g.UnmarshalSession(value)
	if err != nil {
		return nil, nil, err
	}

	err = validateState(rc, sess)
	if err != nil {
		return nil, nil, err
	}

	u, err := g.FetchUser(sess)
	if err == nil {
		return addToSession(u.Provider, u.Email, u.AccessToken, rc, websess, logger)
	}

	_, err = sess.Authorize(g, &params{q: rc.Request.URI().QueryArgs()})
	if err != nil {
		return nil, nil, err
	}

	err = cutil.StoreInSession(prv.ID, sess.Marshal(), rc, websess, logger)
	if err != nil {
		return nil, nil, err
	}

	gu, err := g.FetchUser(sess)
	if err != nil {
		return nil, nil, err
	}

	return addToSession(gu.Provider, gu.Email, gu.AccessToken, rc, websess, logger)
}

func gothFor(rc *fasthttp.RequestCtx, prv *Provider) (goth.Provider, error) {
	proto := string(rc.URI().Scheme())
	host := string(rc.URI().Host())
	if host == "" {
		host = "localhost"
	}
	return prv.Goth(proto, host)
}

func Logout(rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger, prvKeys ...string) error {
	a := getCurrentAuths(websess)
	n := a.Purge(prvKeys...)
	return setCurrentAuths(n, rc, websess, logger)
}
