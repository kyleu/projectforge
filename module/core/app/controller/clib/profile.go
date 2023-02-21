package clib

import (
	"net/url"

	{{{ if .HasModule "oauth" }}}"github.com/pkg/errors"
	{{{ end }}}"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/theme"{{{ if .HasModule "user" }}}
	"{{{ .Package }}}/app/util"{{{ end }}}
	"{{{ .Package }}}/views/vprofile"
)

func Profile(rc *fasthttp.RequestCtx) {
	controller.Act("profile", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(rc, as, ps)
	})
}
{{{ if .HasModule "marketing" }}}
func ProfileSite(rc *fasthttp.RequestCtx) {
	controller.ActSite("profile", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(rc, as, ps)
	})
}
{{{ end }}}
func profileAction(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	ps.Title = "Profile"
	ps.Data = ps.Profile
	thm := as.Themes.Get(ps.Profile.Theme, ps.Logger){{{ if .HasModule "oauth" }}}

	prvs, err := as.Auth.Providers(ps.Logger)
	if err != nil {
		return "", errors.Wrap(err, "can't load providers")
	}{{{ end }}}

	redir := "/"
	ref := string(rc.Request.Header.Peek("Referer"))
	if ref != "" {
		u, err := url.Parse(ref)
		if err == nil && u != nil && u.Path != cutil.DefaultProfilePath {
			redir = u.Path
		}
	}

	page := &vprofile.Profile{Profile: ps.Profile, Theme: thm, {{{ if .HasModule "oauth" }}}Providers: prvs, {{{ end }}}Referrer: redir}
	return controller.Render(rc, as, page, ps, "Profile")
}

func ProfileSave(rc *fasthttp.RequestCtx) {
	controller.Act("profile.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		n := ps.Profile.Clone()

		referrerDefault := frm.GetStringOpt("referrer")
		if referrerDefault == "" {
			referrerDefault = cutil.DefaultProfilePath
		}

		n.Name = frm.GetStringOpt("name")
		n.Mode = frm.GetStringOpt("mode")
		n.Theme = frm.GetStringOpt("theme")
		if n.Theme == theme.ThemeDefault.Key {
			n.Theme = ""
		}{{{ if .HasModule "user" }}}
		if ps.Profile.ID == util.UUIDDefault {
			n.ID = util.UUID()
		} else {
			n.ID = ps.Profile.ID
		}{{{ end }}}

		err = csession.SaveProfile(n, rc, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}{{{ if .HasModule "user" }}}

		curr, _ := as.Services.User.Get(ps.Context, nil, ps.Profile.ID, ps.Logger)
		if curr != nil {
			curr.Name = n.Name
			if curr.Picture == "" {
				curr.Picture = ps.Accounts.Image()
			}
			err = as.Services.User.Update(ps.Context, nil, curr, ps.Logger)
			if err != nil {
				return "", err
			}
		}{{{ end }}}

		return controller.ReturnToReferrer("Saved profile", referrerDefault, rc, ps)
	})
}
