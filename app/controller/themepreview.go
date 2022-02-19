// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/app/lib/telemetry"
	"github.com/kyleu/projectforge/views"
	"github.com/muesli/gamut"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/lib/theme"
	"github.com/kyleu/projectforge/views/vtheme"
)

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
		page := &views.Debug{} // &vtheme.Preview{Theme: thm}
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
