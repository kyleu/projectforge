package controller

import (
	"github.com/kyleu/projectforge/views"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectDetail(ctx *fasthttp.RequestCtx) {
	act("search", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", true)
		if err != nil {
			return "", err
		}

		results := key

		ps.Title = "Search Results"
		ps.Data = results
		return render(ctx, as, &views.Debug{}, ps, "project", key)
	})
}
