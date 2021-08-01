package controller

import (
	"net/url"

	"github.com/go-gem/sessions"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/theme"
	"{{{ .Package }}}/app/user"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vauth"
)

func Profile(ctx *fasthttp.RequestCtx) {
	act("profile", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(ctx, as, ps)
	})
}

func ProfileSite(ctx *fasthttp.RequestCtx) {
	actSite("profile", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(ctx, as, ps)
	})
}

func profileAction(ctx *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	ps.Title = "Profile"
	ps.Data = ps.Profile
	thm := as.Themes.Get(ps.Profile.Theme)

	prvs, err := as.Auth.Providers()
	if err != nil {
		return "", errors.Wrap(err, "can't load providers")
	}

	redir := "/"
	ref := string(ctx.Request.Header.Peek("Referer"))
	if ref != "" {
		u, err := url.Parse(ref)
		if err == nil && u != nil {
			redir = u.Path
		}
	}

	page := &vauth.Profile{Profile: ps.Profile, Theme: thm, Providers: prvs, Referrer: redir}
	return render(ctx, as, page, ps, "Profile")
}

func ProfileSave(ctx *fasthttp.RequestCtx) {
	act("profile.save", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(ctx)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse form")
		}

		n := ps.Profile.Clone()

		n.Name, _ = frm.GetString("name", true)
		n.Mode, _ = frm.GetString("mode", true)
		n.Theme, _ = frm.GetString("theme", true)
		if n.Theme == theme.ThemeDefault.Key {
			n.Theme = ""
		}

		err = user.SaveProfile(n, ctx, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "profile save"
		return returnToReferrer(msg, "/profile", ctx, ps)
	})
}

func loadProfile(session *sessions.Session) (*user.Profile, error) {
	x, ok := session.Values["profile"]
	if !ok {
		return user.DefaultProfile.Clone(), nil
	}
	s, ok := x.(string)
	if !ok {
		return user.DefaultProfile.Clone(), nil
	}
	p := &user.Profile{}
	err := util.FromJSON([]byte(s), p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
