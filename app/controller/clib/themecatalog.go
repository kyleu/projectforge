package clib

import (
	"fmt"
	"net/http"

	"github.com/muesli/gamut"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vtheme"
)

func ThemeColor(w http.ResponseWriter, r *http.Request) {
	controller.Act("theme.color", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		c, err := cutil.PathRichString(r, "color", false)
		if err != nil {
			return "", err
		}
		col := c.ToLower()
		if !col.HasPrefix("#") {
			col = "#" + col
		}
		th := theme.ColorTheme(col.String(), gamut.Hex(col.String()))
		ps.SetTitleAndData(fmt.Sprintf("[%s] Theme", col), th)
		ps.DefaultNavIcon = themeIcon
		return controller.Render(r, as, &vtheme.Edit{Theme: th, Icon: "app"}, ps, "Themes||/theme", col.String())
	})
}

func ThemeColorEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("theme.color.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		color := util.Str(r.URL.Query().Get("color"))
		if color.Empty() {
			return "", errors.New("must provide color in query string")
		}
		if !color.HasPrefix("#") {
			return "", errors.New("provided color must be a hex string")
		}
		t := theme.ColorTheme(color.TrimPrefix("#").String(), gamut.Hex(color.String()))
		ps.SetTitleAndData("Edit theme colors ["+t.Key+"]", t)
		ps.DefaultNavIcon = themeIcon
		page := &vtheme.Edit{Theme: t, Icon: "app", Exists: as.Themes.FileExists(t.Key)}
		return controller.Render(r, as, page, ps, "Themes||/theme", t.Key)
	})
}

func ThemePalette(w http.ResponseWriter, r *http.Request) {
	controller.Act("theme.palette", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pal, err := cutil.PathString(r, "palette", false)
		if err != nil {
			return "", err
		}
		_, span, _ := telemetry.StartSpan(ps.Context, "theme:load", ps.Logger)
		thms, err := theme.PaletteThemes(pal)
		span.Complete()
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Themes", pal), thms)
		if r.URL.Query().Get("t") == "go" {
			ps.Data = util.StringJoin(lo.Map(thms, func(t *theme.Theme, _ int) string {
				return t.ToGo()
			}), util.StringDefaultLinebreak)
			return controller.Render(r, as, &views.Debug{}, ps, "Themes")
		}
		ps.DefaultNavIcon = themeIcon
		return controller.Render(r, as, &vtheme.Add{Palette: pal, Themes: thms}, ps, "Themes||/theme", "Palette")
	})
}

func ThemePaletteEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("theme.palette.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		palette, err := cutil.PathString(r, "palette", false)
		if err != nil {
			return "", err
		}
		key, err := cutil.PathString(r, "theme", false)
		if err != nil {
			return "", err
		}
		if key == theme.Default.Key {
			return controller.FlashAndRedir(false, "Unable to edit default theme", "/theme", ps)
		}
		themes, err := theme.PaletteThemes(palette)
		if err != nil {
			return "", err
		}
		t := themes.Get(key)
		if t == nil {
			return "", errors.Errorf("invalid theme [%s] for palette [%s]", key, palette)
		}
		ps.SetTitleAndData("Edit theme palette ["+t.Key+"]", t)
		ps.DefaultNavIcon = themeIcon
		return controller.Render(r, as, &vtheme.Edit{Theme: t, Icon: "app"}, ps, "Themes||/theme", t.Key)
	})
}
