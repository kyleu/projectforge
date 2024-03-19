package auth

import (
	"net/http"

	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const defaultProfilePath = "/profile"

func BeginAuthHandler(prv *Provider, w http.ResponseWriter, r *http.Request, websess util.ValueMap, logger util.Logger) (string, error) {
	u, err := getAuthURL(prv, w, r, websess, logger)
	if err != nil {
		return "", err
	}
	refer := r.URL.Query().Get("refer")
	if refer != "" && refer != defaultProfilePath {
		_ = csession.StoreInSession(csession.ReferKey, refer, w, websess, logger)
	}
	return u, nil
}

func CompleteUserAuth(prv *Provider, w http.ResponseWriter, r *http.Request, websess util.ValueMap, logger util.Logger) (*user.Account, user.Accounts, error) {
	value, err := csession.GetFromSession(prv.ID, websess)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = removeProviderData(w, websess, logger)
	}()

	g, err := gothFor(r, prv)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create oauth provider")
	}

	sess, err := g.UnmarshalSession(value)
	if err != nil {
		return nil, nil, err
	}

	err = validateState(w, r, sess)
	if err != nil {
		return nil, nil, err
	}

	u, err := g.FetchUser(sess)
	if err == nil {
		return addToSession(u.Provider, u.Email, u.AvatarURL, u.AccessToken, w, websess, logger)
	}

	_, err = sess.Authorize(g, r.URL.Query())
	if err != nil {
		return nil, nil, err
	}

	err = csession.StoreInSession(prv.ID, sess.Marshal(), w, websess, logger)
	if err != nil {
		return nil, nil, err
	}

	gu, err := g.FetchUser(sess)
	if err != nil {
		return nil, nil, err
	}

	return addToSession(gu.Provider, gu.Email, gu.AvatarURL, gu.AccessToken, w, websess, logger)
}

func gothFor(r *http.Request, prv *Provider) (goth.Provider, error) {
	proto := r.URL.Scheme
	host := r.URL.Host
	if host == "" {
		host = "localhost"
	}
	return prv.Goth(proto, host)
}

func Logout(w http.ResponseWriter, r *http.Request, websess util.ValueMap, logger util.Logger, prvKeys ...string) error {
	a := getCurrentAuths(websess)
	n := a.Purge(prvKeys...)
	return setCurrentAuths(n, w, websess, logger)
}
