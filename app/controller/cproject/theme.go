package cproject

import (
	"fmt"
	"net/http"

	"github.com/muesli/gamut"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vproject"
)

func ProjectThemePalette(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.theme.palette", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pal, err := cutil.PathString(r, "palette", false)
		if err != nil {
			return "", err
		}
		prj, err := cutil.PathString(r, "key", false)
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
		if cutil.QueryStringString(ps.URI, "t") == "go" {
			ps.Data = util.StringJoin(lo.Map(x, func(t *theme.Theme, _ int) string {
				return t.ToGo()
			}), util.StringDefaultLinebreak)
			return controller.Render(r, as, &views.Debug{}, ps, "admin", "Themes")
		}
		span.Complete()
		page := &vproject.ThemePalette{Project: prj, Icon: prjIcon, Palette: pal, Themes: x, Title: prjTitle}
		return controller.Render(r, as, page, ps, "admin", "Themes")
	})
}

func ProjectThemeSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.theme.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		pallette, thm, err := themeFromRC(r)
		if err != nil {
			return "", err
		}
		prj, err := cutil.PathString(r, "key", false)
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
		return controller.FlashAndRedir(true, msg, p.WebPath(), ps)
	})
}

func themeFromRC(r *http.Request) (string, *theme.Theme, error) {
	if color, err := cutil.PathString(r, "color", false); err == nil {
		return color, theme.ColorTheme(color, gamut.Hex(color)), nil
	}
	pal, err := cutil.PathString(r, "palette", false)
	if err != nil {
		return "", nil, err
	}
	themeKey, err := cutil.PathString(r, "theme", false)
	if err != nil {
		return "", nil, err
	}
	ret, err := theme.PaletteTheme(pal, themeKey)
	return pal, ret, err
}
