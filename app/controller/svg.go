package controller

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/svg"
	"github.com/kyleu/projectforge/views/vsvg"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func SVGList(ctx *fasthttp.RequestCtx) {
	act("svg.list", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}
		fs := as.Services.Projects.GetFilesystem(prj)

		icons, contents, err := svg.Contents(fs)
		if err != nil {
			return "", errors.Wrap(err, "unable to list project SVGs")
		}

		ps.Title = "SVG Tools"
		ps.Data = icons
		return render(ctx, as, &vsvg.List{Project: prj, Keys: icons, Contents: contents}, ps, "projects", prj.Key, "SVG")
	})
}

func SVGBuild(ctx *fasthttp.RequestCtx) {
	act("svg.build", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}
		fs := as.Services.Projects.GetFilesystem(prj)
		count, err := svg.Build(fs)
		if err != nil {
			return "", err
		}

		msg := fmt.Sprintf("Parsed [%d] SVG files", count)
		return flashAndRedir(true, msg, "/p/"+prj.Key+"/svg", ctx, ps)
	})
}

func SVGAdd(ctx *fasthttp.RequestCtx) {
	act("svg.add", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		qa := ctx.URI().QueryArgs()
		src := string(qa.Peek("src"))
		if src == "" {
			return ersp("must provide [src]")
		}
		tgt := string(qa.Peek("tgt"))
		if tgt == "" {
			tgt = strings.TrimSuffix(src, "-solid")
		}

		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}
		fs := as.Services.Projects.GetFilesystem(prj)

		x, err := svg.AddToProject(fs, src, tgt)
		if err != nil {
			return "", err
		}
		ps.Data = x
		return render(ctx, as, &vsvg.View{SVG: x}, ps, "tools", "svg")
	})
}

func SVGDetail(ctx *fasthttp.RequestCtx) {
	act("svg.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(ctx, as)
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
		return render(ctx, as, &vsvg.View{Project: prj, SVG: x}, ps, "projects", prj.Key, "SVG||/p/"+prj.Key+"/svg", key)
	})
}

func SVGSetApp(ctx *fasthttp.RequestCtx) {
	act("svg.setapp", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(ctx, as)
		if err != nil {
			return "", err
		}
		content, err := svg.Content(fs, key)
		if err != nil {
			return "", errors.Wrap(err, "unable to read SVG ["+key+"]")
		}
		err = svg.SetAppIcon(fs, &svg.SVG{Key: key, Markup: content})
		if err != nil {
			return "", errors.Wrap(err, "unable to set app icon to ["+key+"]")
		}
		msg := "set SVG [" + key + "] as app icon"
		return flashAndRedir(true, msg, "/p/"+prj.Key+"/svg", ctx, ps)
	})
}

func SVGRemove(ctx *fasthttp.RequestCtx) {
	act("svg.remove", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, fs, key, err := prjAndIcon(ctx, as)
		if err != nil {
			return "", err
		}
		if key == "app" {
			return "", errors.New("you can't remove the app icon")
		}
		err = svg.Remove(fs, key)
		if err != nil {
			return "", errors.Wrap(err, "unable to remove SVG ["+key+"]")
		}
		msg := "removed SVG [" + key + "]"
		return flashAndRedir(true, msg, "/p/"+prj.Key+"/svg", ctx, ps)
	})
}

func prjAndIcon(ctx *fasthttp.RequestCtx, as *app.State) (*project.Project, filesystem.FileLoader, string, error) {
	prj, err := getProject(ctx, as)
	if err != nil {
		return nil, nil, "", err
	}
	fs := as.Services.Projects.GetFilesystem(prj)

	key, err := ctxRequiredString(ctx, "icon", false)
	if err != nil {
		return nil, nil, "", err
	}
	return prj, fs, key, nil
}
