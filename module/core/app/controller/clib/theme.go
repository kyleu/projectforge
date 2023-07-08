package clib

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/csession"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/theme"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vtheme"
)

func ThemeList(rc *fasthttp.RequestCtx) {
	controller.Act("theme.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Themes"
		x := as.Themes.All(ps.Logger)
		ps.Data = x
		return controller.Render(rc, as, &vtheme.List{Themes: x}, ps, "admin", "Themes")
	})
}

func ThemeEdit(rc *fasthttp.RequestCtx) {
	controller.Act("theme.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		if key == theme.ThemeDefault.Key {
			return controller.FlashAndRedir(false, "Unable to edit default theme", "/theme", rc, ps)
		}
		var t *theme.Theme
		if key == theme.KeyNew {
			t = theme.ThemeDefault.Clone(key)
		} else {
			t = as.Themes.Get(key, ps.Logger)
		}
		ps.Data = t
		ps.Title = "Edit theme [" + t.Key + "]"
		page := &vtheme.Edit{Theme: t, Icon: "app"}
		return controller.Render(rc, as, page, ps, "admin", "Themes||/admin/theme", t.Key)
	})
}

func ThemeSave(rc *fasthttp.RequestCtx) {
	controller.Act("theme.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		orig := as.Themes.Get(key, ps.Logger)

		newKey, err := frm.GetString("key", false)
		if err != nil {
			return "", err
		}
		if newKey == theme.KeyNew {
			newKey = util.RandomString(12)
		}

		l := orig.Light.Clone().ApplyMap(frm, "light-")
		d := orig.Dark.Clone().ApplyMap(frm, "dark-")

		t := &theme.Theme{Key: newKey, Light: l, Dark: d}

		err = as.Themes.Save(t, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to save theme")
		}

		ps.Profile.Theme = newKey
		err = csession.SaveProfile(ps.Profile, rc, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}

		return controller.ReturnToReferrer("saved changes to theme ["+newKey+"]", "/", rc, ps)
	})
}
