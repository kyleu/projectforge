package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/views/vgraphql"
)

func GraphiQL(rc *fasthttp.RequestCtx) {
	Act("graphiql", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "GraphQL Workspace"
		return Render(rc, as, &vgraphql.Detail{}, ps, "graphql")
	})
}
