package auth

import (
	"net/http"

	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

func getAuthURL(prv *Provider, w http.ResponseWriter, r *http.Request, websess util.ValueMap, logger util.Logger) (string, error) {
	g, err := gothFor(r, prv)
	if err != nil {
		return "", err
	}

	sess, err := g.BeginAuth(getState(r))
	if err != nil {
		return "", err
	}

	u, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	err = csession.StoreInSession(prv.ID, sess.Marshal(), w, websess, logger)
	if err != nil {
		return "", err
	}

	return u, err
}

func getCurrentAuths(websess util.ValueMap) user.Accounts {
	authS, err := csession.GetFromSession(WebAuthKey, websess)
	var ret user.Accounts
	if err == nil && authS != "" {
		ret = user.AccountsFromString(authS)
	}
	return ret
}

func setCurrentAuths(s user.Accounts, w http.ResponseWriter, websess util.ValueMap, logger util.Logger) error {
	s.Sort()
	if len(s) == 0 {
		return csession.RemoveFromSession(WebAuthKey, w, websess, logger)
	}
	return csession.StoreInSession(WebAuthKey, s.String(), w, websess, logger)
}
