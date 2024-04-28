package cproject

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/svg"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vpage"
	"projectforge.dev/projectforge/views/vsvg"
)

const appString, svgLink, svgPath = "app", "SVG||/svg/", "/svg/"

func SVGList(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		icons, contents, err := svg.Contents(pfs, ps.Logger, prj.Modules...)
		if err != nil {
			return "", errors.Wrap(err, "unable to list project SVGs")
		}

		ps.SetTitleAndData("SVG Tools", icons)
		return controller.Render(r, as, &vsvg.List{Project: prj, Keys: icons, Contents: contents}, ps, "projects", prj.Key, "SVG")
	})
}

func SVGDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(r, as)
		if err != nil {
			return "", err
		}
		content, err := svg.Content(fs, key, prj.Modules...)
		if err != nil {
			return "", errors.Wrap(err, "unable to read SVG ["+key+"]")
		}
		x := &svg.SVG{Key: key, Markup: content}
		ps.SetTitleAndData("SVG ["+key+"]", x)
		return controller.Render(r, as, &vsvg.View{Project: prj, SVG: x}, ps, "projects", prj.Key, svgLink+prj.Key, key)
	})
}

func SVGBuild(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.build", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		count, err := svg.Build(pfs, ps.Logger, prj)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("Parsed [%d] SVG files", count)
		return controller.FlashAndRedir(true, msg, svgPath+prj.Key, ps)
	})
}

func SVGAdd(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.add", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		qa := r.URL.Query()
		src := strings.TrimSpace(qa.Get("src"))
		if src == "" {
			return controller.ERsp("must provide [src]")
		}
		tgt := qa.Get("tgt")
		if tgt == "" {
			tgt = strings.TrimSuffix(src, "-solid")
		}
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		x, err := svg.AddToProject(pfs, src, tgt, prj.Modules...)
		if err != nil {
			return "", err
		}
		_, err = svg.Build(pfs, ps.Logger, prj)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Added SVG ["+x.Key+"]", x)
		return controller.Render(r, as, &vsvg.View{Project: prj, SVG: x}, ps, "projects", prj.Key, svgLink+prj.Key, x.Key)
	})
}

func SVGSetApp(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.set.app", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(r, as)
		if err != nil {
			return "", err
		}
		if r.URL.Query().Get("hasloaded") != util.BoolTrue {
			cutil.URLAddQuery(r.URL, "hasloaded", util.BoolTrue)
			page := &vpage.Load{URL: r.URL.String(), Title: "Generating icons"}
			return controller.Render(r, as, page, ps, "projects", prj.Key, svgLink+prj.Key, "App Icon")
		}
		content, err := svg.Content(fs, key, prj.Modules...)
		if err != nil {
			return "", errors.Wrap(err, "unable to read app SVG ["+key+"]")
		}
		prj.Icon = key
		err = as.Services.Projects.Save(prj, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to set project icon ["+key+"]")
		}
		err = svg.SetAppIcon(ps.Context, prj, fs, &svg.SVG{Key: key, Markup: content}, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to set app icon to ["+key+"]")
		}
		_, err = svg.Build(fs, ps.Logger, prj)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "set SVG ["+key+"] as app icon", svgPath+prj.Key, ps)
	})
}

func SVGRefreshApp(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.refresh.app", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		if r.URL.Query().Get("hasloaded") != util.BoolTrue {
			cutil.URLAddQuery(r.URL, "hasloaded", util.BoolTrue)
			page := &vpage.Load{URL: r.URL.String(), Title: "Generating app icons"}
			return controller.Render(r, as, page, ps, "projects", prj.Key, svgLink+prj.Key, "Refresh App Icon")
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = svg.RefreshAppIcon(ps.Context, prj, pfs, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to refresh app icon")
		}
		return controller.FlashAndRedir(true, "refreshed app icon", svgPath+prj.Key, ps)
	})
}

func SVGRemove(w http.ResponseWriter, r *http.Request) {
	controller.Act("svg.remove", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(r, as)
		if err != nil {
			return "", err
		}
		if key == appString {
			return "", errors.New("you can't remove the app icon")
		}
		err = svg.Remove(fs, key, ps.Logger, prj.Modules...)
		if err != nil {
			return "", errors.Wrap(err, "unable to remove SVG ["+key+"]")
		}
		_, err = svg.Build(fs, ps.Logger, prj)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "removed SVG ["+key+"]", svgPath+prj.Key, ps)
	})
}

func prjAndIcon(r *http.Request, as *app.State) (*project.Project, filesystem.FileLoader, string, error) {
	prj, err := getProject(r, as)
	if err != nil {
		return nil, nil, "", err
	}
	pfs, err := as.Services.Projects.GetFilesystem(prj)
	if err != nil {
		return nil, nil, "", err
	}

	key, err := cutil.PathString(r, "icon", false)
	if err != nil {
		return nil, nil, "", err
	}
	return prj, pfs, key, nil
}
