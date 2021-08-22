package controller

import (
	"strings"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/views/vsearch"
)

func Search(rc *fasthttp.RequestCtx) {
	act("search", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		q := string(rc.URI().QueryArgs().Peek("q"))
		q = strings.TrimSpace(q)
		results := []string{"a", "b", "c"}
		ps.Title = "Search Results"
		ps.Data = results
		return render(rc, as, &vsearch.Results{Q: q, Results: results}, ps, "Search")
	})
}
