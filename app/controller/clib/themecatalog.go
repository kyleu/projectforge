// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"fmt"
	"strings"

	"github.com/muesli/gamut"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vtheme"
)

func ThemeColor(rc *fasthttp.RequestCtx) {
	controller.Act("theme.color", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		col, err := cutil.RCRequiredString(rc, "color", false)
		if err != nil {
			return "", err
		}
		col = strings.ToLower(col)
		if !strings.HasPrefix(col, "#") {
			col = "#" + col
		}
		th := theme.ColorTheme(col, gamut.Hex(col))
		ps.SetTitleAndData(fmt.Sprintf("[%s] Theme", col), th)
		ps.DefaultNavIcon = themeIcon
		return controller.Render(rc, as, &vtheme.Edit{Theme: th, Icon: "app"}, ps, "Themes||/theme", col)
	})
}

func ThemeColorEdit(rc *fasthttp.RequestCtx) {
	controller.Act("theme.color.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		color := string(rc.URI().QueryArgs().Peek("color"))
		if color == "" {
			return "", errors.New("must provide color in query string")
		}
		if !strings.HasPrefix(color, "#") {
			return "", errors.New("provided color must be a hex string")
		}
		t := theme.ColorTheme(strings.TrimPrefix(color, "#"), gamut.Hex(color))
		ps.SetTitleAndData("Edit theme ["+t.Key+"]", t)
		ps.DefaultNavIcon = themeIcon
		page := &vtheme.Edit{Theme: t, Icon: "app", Exists: as.Themes.FileExists(t.Key)}
		return controller.Render(rc, as, page, ps, "Themes||/theme", t.Key)
	})
}

func ThemePalette(rc *fasthttp.RequestCtx) {
	controller.Act("theme.palette", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pal, err := cutil.RCRequiredString(rc, "palette", false)
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
		if string(rc.URI().QueryArgs().Peek("t")) == "go" {
			ps.Data = strings.Join(lo.Map(thms, func(t *theme.Theme, _ int) string {
				return t.ToGo()
			}), util.StringDefaultLinebreak)
			return controller.Render(rc, as, &views.Debug{}, ps, "Themes")
		}
		ps.DefaultNavIcon = themeIcon
		return controller.Render(rc, as, &vtheme.Add{Palette: pal, Themes: thms}, ps, "Themes||/theme", "Palette")
	})
}

func ThemePaletteEdit(rc *fasthttp.RequestCtx) {
	controller.Act("theme.palette.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		palette, err := cutil.RCRequiredString(rc, "palette", false)
		if err != nil {
			return "", err
		}
		key, err := cutil.RCRequiredString(rc, "theme", false)
		if err != nil {
			return "", err
		}
		if key == theme.Default.Key {
			return controller.FlashAndRedir(false, "Unable to edit default theme", "/theme", rc, ps)
		}
		themes, err := theme.PaletteThemes(palette)
		if err != nil {
			return "", err
		}
		t := themes.Get(key)
		if t == nil {
			return "", errors.Errorf("invalid theme [%s] for palette [%s]", key, palette)
		}
		ps.SetTitleAndData("Edit theme ["+t.Key+"]", t)
		ps.DefaultNavIcon = themeIcon
		return controller.Render(rc, as, &vtheme.Edit{Theme: t, Icon: "app"}, ps, "Themes||/theme", t.Key)
	})
}
