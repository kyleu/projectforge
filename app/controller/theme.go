// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vtheme"
)

func ThemeList(rc *fasthttp.RequestCtx) {
	act("theme.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Themes"
		x := as.Themes.All(ps.Logger)
		ps.Data = x
		return render(rc, as, &vtheme.List{Themes: x}, ps, "admin", "Themes")
	})
}

func ThemeEdit(rc *fasthttp.RequestCtx) {
	act("theme.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		var t *theme.Theme
		if key == theme.KeyNew {
			t = theme.ThemeDefault.Clone(key)
		} else {
			t = as.Themes.Get(key, ps.Logger)
		}
		if t == nil {
			return "", errors.Wrap(err, "no theme found with key ["+key+"]")
		}
		ps.Data = t
		ps.Title = "Edit theme [" + t.Key + "]"
		page := &vtheme.Edit{Theme: t}
		return render(rc, as, page, ps, "admin", "Themes||/admin/theme", t.Key)
	})
}

func ThemeSave(rc *fasthttp.RequestCtx) {
	act("theme.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
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
		err = cutil.SaveProfile(ps.Profile, rc, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}

		return returnToReferrer("saved changes to theme ["+newKey+"]", "/", rc, ps)
	})
}
