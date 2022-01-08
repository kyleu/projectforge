package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/lib/auth"
)

const signinMsg = "signed in using %s as [%s]"

func AuthDetail(rc *fasthttp.RequestCtx) {
	act("auth.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prv, err := getProvider(as, rc)
		if err != nil {
			return "", err
		}
		u, _, err := auth.CompleteUserAuth(prv, rc, ps.Session, ps.Logger)
		if err == nil {
			msg := fmt.Sprintf(signinMsg, auth.AvailableProviderNames[prv.ID], u.Email)
			return returnToReferrer(msg, defaultProfilePath, rc, ps)
		}
		return auth.BeginAuthHandler(prv, rc, ps.Session, ps.Logger)
	})
}

func AuthCallback(rc *fasthttp.RequestCtx) {
	act("auth.callback", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prv, err := getProvider(as, rc)
		if err != nil {
			return "", err
		}
		u, _, err := auth.CompleteUserAuth(prv, rc, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf(signinMsg, auth.AvailableProviderNames[prv.ID], u.Email)
		return returnToReferrer(msg, defaultProfilePath, rc, ps)
	})
}

func AuthLogout(rc *fasthttp.RequestCtx) {
	act("auth.logout", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		err = auth.Logout(rc, ps.Session, ps.Logger, key)
		if err != nil {
			return "", err
		}

		return ps.ProfilePath, nil
	})
}

func getProvider(as *app.State, rc *fasthttp.RequestCtx) (*auth.Provider, error) {
	key, err := RCRequiredString(rc, "key", false)
	if err != nil {
		return nil, err
	}
	prvs, err := as.Auth.Providers()
	if err != nil {
		return nil, errors.Wrap(err, "can't load providers")
	}
	prv := prvs.Get(key)
	if prv == nil {
		return nil, errors.Errorf("no provider available with id [%s]", key)
	}
	return prv, nil
}
