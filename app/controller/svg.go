package controller

import (
	"fmt"
	"strings"

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

		icons, err := svg.List(fs)
		if err != nil {
			return "", errors.Wrap(err, "unable to list project SVGs")
		}

		ps.Title = "SVG Tools"
		ps.Data = icons
		return render(ctx, as, &vsvg.List{Project: prj}, ps, "projects", prj.Key, "SVG Tools")
	})
}

func SVGBuild(ctx *fasthttp.RequestCtx) {
	act("svg.build", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(ctx, as)
		if err != nil {
			return "", err
		}
		fs := as.Services.Projects.GetFilesystem(prj)
		count, err := svg.Run(fs, "client/src/svg", "app/util/svg.go")
		if err != nil {
			return "", err
		}

		msg := fmt.Sprintf("Parsed [%d] SVG files", count)
		return flashAndRedir(true, msg, "", ctx, ps)
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
