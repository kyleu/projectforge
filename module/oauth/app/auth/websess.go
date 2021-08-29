package auth

import (
	"github.com/go-gem/sessions"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/user"
)

const WebSessKey = "auth"

func addToSession(provider string, email string, rc *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) (*user.Account, user.Accounts, error) {
	ret := getCurrentAuths(websess)
	s := &user.Account{Provider: provider, Email: email}
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

func removeProviderData(rc *fasthttp.RequestCtx, websess *sessions.Session, logger *zap.SugaredLogger) error {
	dirty := false
	for k := range websess.Values {
		s, ok := k.(string)
		if !ok {
			logger.Error("unable to parse session key [%s] of type [%T]", k, k)
		}
		if isProvider(s) {
			logger.Debug("removing auth info for provider [" + s + "]")
			dirty = true
			delete(websess.Values, k)
		}
	}
	if dirty {
		return cutil.SaveSession(rc, websess, logger)
	}
	return nil
}

func isProvider(k string) bool {
	for _, x := range AvailableProviderKeys {
		if x == k {
			return true
		}
	}
	return false
}
