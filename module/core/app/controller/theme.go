package controller

import (
	"$PF_PACKAGE$/app/theme"
	"$PF_PACKAGE$/app/user"
	"$PF_PACKAGE$/app/util"
	"$PF_PACKAGE$/views/vtheme"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app/controller/cutil"

	"$PF_PACKAGE$/app"
)

func ThemeList(ctx *fasthttp.RequestCtx) {
	act("theme.list", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Themes"
		x := as.Themes.All()
		ps.Data = x
		return render(ctx, as, &vtheme.List{Themes: x}, ps, "settings", "Themes||/theme")
	})
}

func ThemeEdit(ctx *fasthttp.RequestCtx) {
	act("theme.edit", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", false)
		if err != nil {
			return "", err
		}
		var t *theme.Theme
		if key == theme.KeyNew {
			t = theme.ThemeDefault.Clone(key)
		} else {
			t = as.Themes.Get(key)
		}
		if t == nil {
			return "", errors.Wrap(err, "no theme found with key ["+key+"]")
		}
		ps.Data = t
		ps.Title = "Edit theme [" + t.Key + "]"
		page := &vtheme.Edit{Theme: t}
		return render(ctx, as, page, ps, "settings", "Themes||/theme", t.Key)
	})
}

func ThemeSave(ctx *fasthttp.RequestCtx) {
	act("theme.save", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", false)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(ctx)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse form")
		}

		orig := as.Themes.Get(key)

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

		err = as.Themes.Save(t)
		if err != nil {
			return "", errors.Wrap(err, "unable to save theme")
		}

		ps.Profile.Theme = newKey
		err = user.SaveProfile(ps.Profile, ctx, ps.Session, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "saved changes to theme [" + newKey + "]"
		return returnToReferrer(msg, "/", ctx, ps)
	})
}
