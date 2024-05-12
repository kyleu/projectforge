package auth

import (
	"net/http"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const WebAuthKey = "auth"

func addToSession(
	provider string, name string, email string, picture string, token string, w http.ResponseWriter, websess util.ValueMap, logger util.Logger,
) (string, *user.Account, user.Accounts, error) {
	ret := getCurrentAuths(websess)
	s := &user.Account{Provider: provider, Email: email, Picture: picture, Token: token}
	for _, x := range ret {
		if x.Matches(s) {
			return name, s, ret, nil
		}
	}
	ret = append(ret, s)
	err := setCurrentAuths(ret, w, websess, logger)
	if err != nil {
		return name, nil, nil, err
	}
	return name, s, ret, nil
}

func removeProviderData(w http.ResponseWriter, websess util.ValueMap, logger util.Logger) error {
	dirty := false
	lo.ForEach(websess.Keys(), func(s string, _ int) {
		if isProvider(s) {
			logger.Debug("removing auth info for provider [" + s + "]")
			dirty = true
			delete(websess, s)
		}
	})
	if dirty {
		return csession.SaveSession(w, websess, logger)
	}
	return nil
}

func isProvider(k string) bool {
	return lo.Contains(AvailableProviderKeys, k)
}
