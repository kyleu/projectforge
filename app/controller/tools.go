package controller

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/svg"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vsvg"
	"github.com/kyleu/projectforge/views/vtools"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ToolList(ctx *fasthttp.RequestCtx) {
	act("tool.list", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Tools"
		ps.Data = util.ValueMap{
			"svg": "/tools/svg",
		}
		return render(ctx, as, &vtools.List{}, ps, "tools")
	})
}

func SVGList(ctx *fasthttp.RequestCtx) {
	act("svg.list", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "SVG Tools"
		ps.Data = util.ValueMap{"available": util.SVGIconKeys}
		return render(ctx, as, &vsvg.List{}, ps, "tools", "svg")
	})
}

func SVGBuild(ctx *fasthttp.RequestCtx) {
	act("svg.build", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		count, err := svg.Run("client/src/svg", "app/util/svg.go")
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
		x, err := svg.Save(src, tgt)
		if err != nil {
			return "", err
		}
		ps.Data = x
		return render(ctx, as, &vsvg.View{SVG: x}, ps, "tools", "svg")
	})
}
