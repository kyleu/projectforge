package cproject

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

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

const appString = "app"

func SVGList(rc *fasthttp.RequestCtx) {
	controller.Act("svg.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		icons, contents, err := svg.Contents(pfs, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to list project SVGs")
		}

		ps.Title = "SVG Tools"
		ps.Data = icons
		return controller.Render(rc, as, &vsvg.List{Project: prj, Keys: icons, Contents: contents}, ps, "projects", prj.Key, "SVG")
	})
}

func SVGDetail(rc *fasthttp.RequestCtx) {
	controller.Act("svg.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(rc, as)
		if err != nil {
			return "", err
		}
		content, err := svg.Content(fs, key)
		if err != nil {
			return "", errors.Wrap(err, "unable to read SVG ["+key+"]")
		}
		x := &svg.SVG{Key: key, Markup: content}
		ps.Title = "SVG [" + key + "]"
		ps.Data = x
		return controller.Render(rc, as, &vsvg.View{Project: prj, SVG: x}, ps, "projects", prj.Key, "SVG||/svg/"+prj.Key, key)
	})
}

func SVGBuild(rc *fasthttp.RequestCtx) {
	controller.Act("svg.build", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		count, err := svg.Build(pfs, ps.Logger)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("Parsed [%d] SVG files", count)
		return controller.FlashAndRedir(true, msg, "/svg/"+prj.Key, rc, ps)
	})
}

func SVGAdd(rc *fasthttp.RequestCtx) {
	controller.Act("svg.add", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		qa := rc.URI().QueryArgs()
		src := strings.TrimSpace(string(qa.Peek("src")))
		if src == "" {
			return controller.ERsp("must provide [src]")
		}
		tgt := string(qa.Peek("tgt"))
		if tgt == "" {
			tgt = strings.TrimSuffix(src, "-solid")
		}
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		x, err := svg.AddToProject(pfs, src, tgt)
		if err != nil {
			return "", err
		}
		_, err = svg.Build(pfs, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Title = "SVG [" + x.Key + "]"
		ps.Data = x
		return controller.Render(rc, as, &vsvg.View{Project: prj, SVG: x}, ps, "projects", prj.Key, "SVG||/svg/"+prj.Key, x.Key)
	})
}

func SVGSetApp(rc *fasthttp.RequestCtx) {
	controller.Act("svg.set.app", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(rc, as)
		if err != nil {
			return "", err
		}
		if string(rc.URI().QueryArgs().Peek("hasloaded")) != util.BoolTrue {
			rc.URI().QueryArgs().Set("hasloaded", util.BoolTrue)
			page := &vpage.Load{URL: rc.URI().String(), Title: "Generating icons"}
			return controller.Render(rc, as, page, ps, "projects", prj.Key, "SVG||/svg/"+prj.Key, "App Icon")
		}
		content, err := svg.Content(fs, key)
		if err != nil {
			return "", errors.Wrap(err, "unable to read SVG ["+key+"]")
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
		_, err = svg.Build(fs, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "set SVG ["+key+"] as app icon", "/svg/"+prj.Key, rc, ps)
	})
}

func SVGRefreshApp(rc *fasthttp.RequestCtx) {
	controller.Act("svg.refresh.app", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		if string(rc.URI().QueryArgs().Peek("hasloaded")) != util.BoolTrue {
			rc.URI().QueryArgs().Set("hasloaded", util.BoolTrue)
			page := &vpage.Load{URL: rc.URI().String(), Title: "Generating app icons"}
			return controller.Render(rc, as, page, ps, "projects", prj.Key, "SVG||/svg/"+prj.Key, "Refresh App Icon")
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = svg.RefreshAppIcon(ps.Context, prj, pfs, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to refresh app icon")
		}
		return controller.FlashAndRedir(true, "refreshed app icon", "/svg/"+prj.Key, rc, ps)
	})
}

func SVGRemove(rc *fasthttp.RequestCtx) {
	controller.Act("svg.remove", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(rc, as)
		if err != nil {
			return "", err
		}
		if key == appString {
			return "", errors.New("you can't remove the app icon")
		}
		err = svg.Remove(fs, key, ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to remove SVG ["+key+"]")
		}
		_, err = svg.Build(fs, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "removed SVG ["+key+"]", "/svg/"+prj.Key, rc, ps)
	})
}

func prjAndIcon(rc *fasthttp.RequestCtx, as *app.State) (*project.Project, filesystem.FileLoader, string, error) {
	prj, err := getProject(rc, as)
	if err != nil {
		return nil, nil, "", err
	}
	pfs, err := as.Services.Projects.GetFilesystem(prj)
	if err != nil {
		return nil, nil, "", err
	}

	key, err := cutil.RCRequiredString(rc, "icon", false)
	if err != nil {
		return nil, nil, "", err
	}
	return prj, pfs, key, nil
}
