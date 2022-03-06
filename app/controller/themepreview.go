package controller

import (
	"fmt"
	"strings"

	"github.com/muesli/gamut"
	"github.com/valyala/fasthttp"
	"projectforge.dev/app/lib/telemetry"

	"projectforge.dev/app"
	"projectforge.dev/app/controller/cutil"
	"projectforge.dev/app/lib/theme"
	"projectforge.dev/views/vtheme"
)

func ThemeColor(rc *fasthttp.RequestCtx) {
	act("theme.color", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		col, err := RCRequiredString(rc, "color", false)
		if err != nil {
			return "", err
		}
		col = strings.ToLower(col)
		if !strings.HasPrefix(col, "#") {
			col = "#" + col
		}
		ps.Title = fmt.Sprintf("[%s] Theme", col)
		x := theme.ColorTheme(col, gamut.Hex(col))
		ps.Data = x
		return render(rc, as, &vtheme.Edit{Theme: x}, ps, "admin", "Themes", col)
	})
}

func ThemePalette(rc *fasthttp.RequestCtx) {
	act("theme.palette", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pal, err := RCRequiredString(rc, "palette", false)
		if err != nil {
			return "", err
		}
		ps.Title = fmt.Sprintf("[%s] Themes", pal)
		_, span, _ := telemetry.StartSpan(ps.Context, "theme:load", ps.Logger)
		x, err := theme.PaletteThemes(pal)
		if err != nil {
			return "", err
		}
		span.Complete()
		ps.Data = x
		return render(rc, as, &vtheme.Add{Palette: pal, Themes: x}, ps, "admin", "Themes")
	})
}

func ThemePreview(rc *fasthttp.RequestCtx) {
	act("theme.preview", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		thm, err := themeFromRC(rc)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("Preview [%s]", thm.Key)

		ps.Data = thm
		page := &vtheme.Edit{Theme: thm}
		return render(rc, as, page, ps, "admin", "Themes||/admin/theme", thm.Key)
	})
}

func themeFromRC(rc *fasthttp.RequestCtx) (*theme.Theme, error) {
	color, err := RCRequiredString(rc, "color", false)
	if err != nil {
		pal, err := RCRequiredString(rc, "palette", false)
		if err != nil {
			return nil, err
		}
		key, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return nil, err
		}
		return theme.PaletteTheme(pal, key)
	}
	return theme.ColorTheme(color, gamut.Hex(color)), nil
}
