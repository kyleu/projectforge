// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"net/http"
	"net/url"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/views/vprofile"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	controller.Act("profile", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(w, r, as, ps)
	})
}

func ProfileSite(w http.ResponseWriter, r *http.Request) {
	controller.ActSite("profile", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		return profileAction(w, r, as, ps)
	})
}

func profileAction(w http.ResponseWriter, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	ps.SetTitleAndData("Profile", ps.Profile)
	thm := as.Themes.Get(ps.Profile.Theme, ps.Logger)

	redir := "/"
	ref := r.Header.Get("Referer")
	if ref != "" {
		u, err := url.Parse(ref)
		if err == nil && u != nil && u.Path != cutil.DefaultProfilePath {
			redir = u.Path
		}
	}
	ps.DefaultNavIcon = "profile"
	page := &vprofile.Profile{Profile: ps.Profile, Theme: thm, Referrer: redir}
	return controller.Render(r, as, page, ps, "Profile")
}

func ProfileSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("profile.save", w, r, func(_ *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		if n.Theme == theme.Default.Key {
			n.Theme = ""
		}

		err = csession.SaveProfile(n, w, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}

		return controller.ReturnToReferrer("Saved profile", referrerDefault, ps)
	})
}
