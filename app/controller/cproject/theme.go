package cproject

import (
	"fmt"
	"strings"

	"github.com/muesli/gamut"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectThemePalette(rc *fasthttp.RequestCtx) {
	controller.Act("project.theme.palette", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pal, err := cutil.RCRequiredString(rc, "palette", false)
		if err != nil {
			return "", err
		}
		prj, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		curr, err := as.Services.Projects.Get(prj)
		if err != nil {
			return "", err
		}
		prjTitle := curr.Title()
		prjIcon := curr.IconSafe()
		_, span, _ := telemetry.StartSpan(ps.Context, "theme:load", ps.Logger)
		x, err := theme.PaletteThemes(pal)
		span.Complete()
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Themes", pal), pal)
		if string(rc.URI().QueryArgs().Peek("t")) == "go" {
			ps.Data = strings.Join(lo.Map(x, func(t *theme.Theme, _ int) string {
				return t.ToGo()
			}), util.StringDefaultLinebreak)
			return controller.Render(rc, as, &views.Debug{}, ps, "admin", "Themes")
		}
		span.Complete()
		page := &vproject.ThemePalette{Project: prj, Icon: prjIcon, Palette: pal, Themes: x, Title: prjTitle}
		return controller.Render(rc, as, page, ps, "admin", "Themes")
	})
}

func ProjectThemeSave(rc *fasthttp.RequestCtx) {
	controller.Act("project.theme.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pallette, thm, err := themeFromRC(rc)
		if err != nil {
			return "", err
		}
		prj, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
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
		return controller.FlashAndRedir(true, msg, "/p/"+p.Key, rc, ps)
	})
}

func themeFromRC(rc *fasthttp.RequestCtx) (string, *theme.Theme, error) {
	if color, err := cutil.RCRequiredString(rc, "color", false); err == nil {
		return color, theme.ColorTheme(color, gamut.Hex(color)), nil
	}
	pal, err := cutil.RCRequiredString(rc, "palette", false)
	if err != nil {
		return "", nil, err
	}
	themeKey, err := cutil.RCRequiredString(rc, "theme", false)
	if err != nil {
		return "", nil, err
	}
	ret, err := theme.PaletteTheme(pal, themeKey)
	return pal, ret, err
}
