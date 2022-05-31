package controller

import (
	"fmt"
	"strings"

	"github.com/muesli/gamut"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vtheme"
)

func ThemeColor(rc *fasthttp.RequestCtx) {
	act("theme.color", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		col, err := cutil.RCRequiredString(rc, "color", false)
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
		pal, err := cutil.RCRequiredString(rc, "palette", false)
		if err != nil {
			return "", err
		}
		prj := string(rc.URI().QueryArgs().Peek("project"))
		prjTitle := util.AppName
		if prj != "" {
			curr, _ := as.Services.Projects.Get(prj)
			if curr != nil {
				prjTitle = curr.Title()
			}
		}
		ps.Title = fmt.Sprintf("[%s] Themes", pal)
		_, span, _ := telemetry.StartSpan(ps.Context, "theme:load", ps.Logger)
		x, err := theme.PaletteThemes(pal)
		if err != nil {
			return "", err
		}
		span.Complete()
		ps.Data = x
		return render(rc, as, &vtheme.Add{Project: prj, Palette: pal, Themes: x, Title: prjTitle}, ps, "admin", "Themes")
	})
}

func ThemePreview(rc *fasthttp.RequestCtx) {
	act("theme.preview", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pallette, thm, err := themeFromRC(rc)
		if err != nil {
			return "", err
		}

		if prj := string(rc.URI().QueryArgs().Peek("project")); prj != "" {
			p, err := as.Services.Projects.Get(prj)
			if err != nil {
				return "", err
			}
			p.Theme = thm
			err = as.Services.Projects.Save(p, ps.Logger)
			if err != nil {
				return "", err
			}
			msg := fmt.Sprintf("set theme to [%s:%s]", pallette, thm.Key)
			return flashAndRedir(true, msg, "/p/"+p.Key, rc, ps)
		}

		ps.Title = fmt.Sprintf("Preview [%s]", thm.Key)
		ps.Data = thm
		page := &vtheme.Edit{Theme: thm}
		return render(rc, as, page, ps, "admin", "Themes||/admin/theme", thm.Key)
	})
}

func themeFromRC(rc *fasthttp.RequestCtx) (string, *theme.Theme, error) {
	color, err := cutil.RCRequiredString(rc, "color", false)
	if err != nil {
		pal, err := cutil.RCRequiredString(rc, "palette", false)
		if err != nil {
			return "", nil, err
		}
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", nil, err
		}
		ret, err := theme.PaletteTheme(pal, key)
		return pal, ret, err
	}
	return color, theme.ColorTheme(color, gamut.Hex(color)), nil
}
