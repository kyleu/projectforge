package auth

import (
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const WebAuthKey = "auth"

func addToSession(
	provider string, email string, picture string, token string, rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger,
) (*user.Account, user.Accounts, error) {
	ret := getCurrentAuths(websess)
	s := &user.Account{Provider: provider, Email: email, Picture: picture, Token: token}
	for _, x := range ret {
		if x.Provider == s.Provider && x.Email == s.Email {
			return s, ret, nil
		}
	}
	ret = append(ret, s)
	err := setCurrentAuths(ret, rc, websess, logger)
	if err != nil {
		return nil, nil, err
	}
	return s, ret, nil
}

func removeProviderData(rc *fasthttp.RequestCtx, websess util.ValueMap, logger util.Logger) error {
	dirty := false
	lo.ForEach(websess.Keys(), func(s string, _ int) {
		if isProvider(s) {
			logger.Debug("removing auth info for provider [" + s + "]")
			dirty = true
			delete(websess, s)
		}
	})
	if dirty {
		return csession.SaveSession(rc, websess, logger)
	}
	return nil
}

func isProvider(k string) bool {
	return lo.Contains(AvailableProviderKeys, k)
}
