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
	"{{{ .Package }}}/app/web"
	"{{{ .Package }}}/views/vprofile"
)

func Profile(ctx *fasthttp.RequestCtx) {
	act("profile", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(ctx, as, ps)
	})
}
{{{ if .HasModule "marketing" }}}
func ProfileSite(ctx *fasthttp.RequestCtx) {
	actSite("profile", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(ctx, as, ps)
	})
}
{{{ end }}}
func profileAction(ctx *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	ps.Title = "Profile"
	ps.Data = ps.Profile
	thm := as.Themes.Get(ps.Profile.Theme)
	{{{ if .HasModule "oauth" }}}
	prvs, err := as.Auth.Providers()
	if err != nil {
		return "", errors.Wrap(err, "can't load providers")
	}
	{{{ end }}}
	redir := "/"
	ref := string(ctx.Request.Header.Peek("Referer"))
	if ref != "" {
		u, err := url.Parse(ref)
		if err == nil && u != nil {
			redir = u.Path
		}
	}

	page := &vprofile.Profile{Profile: ps.Profile, Theme: thm, {{{ if .HasModule "oauth" }}}Providers: prvs, {{{ end }}}Referrer: redir}
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

func returnToReferrer(msg string, dflt string, ctx *fasthttp.RequestCtx, ps *cutil.PageState) (string, error) {
	refer := ""
	referX, ok := ps.Session.Values[web.ReferKey]
	if ok {
		refer, ok = referX.(string)
		if ok {
			_ = web.RemoveFromSession(web.ReferKey, ctx, ps.Session, ps.Logger)
		}
	}
	if refer == "" {
		refer = dflt
	}
	return flashAndRedir(true, msg, refer, ctx, ps)
}
