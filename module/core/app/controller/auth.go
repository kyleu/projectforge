package controller

import (
	"fmt"

	"$PF_PACKAGE$/app/auth"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app/controller/cutil"

	"$PF_PACKAGE$/app"
)

const signinMsg = "signed in to %s as [%s]"

func AuthDetail(ctx *fasthttp.RequestCtx) {
	act("auth.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prv, err := getProvider(as, ctx)
		if err != nil {
			return "", err
		}
		u, _, err := auth.CompleteUserAuth(prv, ctx, ps.Session, ps.Logger)
		if err == nil {
			msg := fmt.Sprintf(signinMsg, auth.AvailableProviderNames[prv.ID], u.Email)
			return returnToReferrer(msg, "/profile", ctx, ps)
		}
		return auth.BeginAuthHandler(prv, ctx, ps.Session, ps.Logger)
	})
}

func AuthCallback(ctx *fasthttp.RequestCtx) {
	act("auth.callback", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prv, err := getProvider(as, ctx)
		if err != nil {
			return "", err
		}
		u, _, err := auth.CompleteUserAuth(prv, ctx, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf(signinMsg, auth.AvailableProviderNames[prv.ID], u.Email)
		return returnToReferrer(msg, "/profile", ctx, ps)
	})
}

func AuthLogout(ctx *fasthttp.RequestCtx) {
	act("auth.logout", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", false)
		if err != nil {
			return "", err
		}
		err = auth.Logout(ctx, ps.Session, ps.Logger, key)
		if err != nil {
			return "", err
		}

		return ps.ProfilePath, nil
	})
}

func returnToReferrer(msg string, dflt string, ctx *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	refer := ""
	referX, ok := ps.Session.Values[auth.ReferKey]
	if ok {
		refer, ok = referX.(string)
		if ok {
			_ = auth.RemoveFromSession(auth.ReferKey, ctx, ps.Session, ps.Logger)
		}
	}
	if refer == "" {
		refer = dflt
	}
	return flashAndRedir(true, msg, refer, ctx, ps)
}

func getProvider(as *app.State, ctx *fasthttp.RequestCtx) (*auth.Provider, error) {
	key, err := ctxRequiredString(ctx, "key", false)
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
