package controller

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vgraphql"
)

func GraphQLIndex(rc *fasthttp.RequestCtx) {
	Act("graphql.index", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "GraphQL List"
		keys := as.GraphQL.Keys()
		if len(keys) == 1 {
			return "/graphql/" + keys[0], nil
		}
		counts := lo.Map(keys, func(key string, _ int) int {
			return as.GraphQL.ExecCount(key)
		})
		ps.Data = keys
		return Render(rc, as, &vgraphql.List{Keys: keys, Counts: counts}, ps, "graphql")
	})
}

func GraphQLDetail(rc *fasthttp.RequestCtx) {
	Act("graphql.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		titles := as.GraphQL.Titles()
		if len(titles) > 1 {
			ps.Title = fmt.Sprintf("[%s] GraphQL Workspace", titles[key])
		} else {
			ps.Title = "GraphQL Workspace"
		}
		ps.Data = key
		bc := []string{"graphql"}
		if len(titles) > 1 {
			bc = append(bc, key)
		}
		return Render(rc, as, &vgraphql.Detail{Key: key}, ps, bc...)
	})
}

func GraphQLRun(rc *fasthttp.RequestCtx) {
	Act("graphql.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}
		q := frm.GetStringOpt("query")
		op := frm.GetStringOpt("operationName")
		v := frm.GetStringOpt("variables")
		vars := util.ValueMap{}
		if v != "" {
			_ = util.FromJSON([]byte(v), &vars)
		}
		rsp, err := as.GraphQL.Exec(ps.Context, key, q, op, vars, ps.Logger)
		if err != nil {
			return "", err
		}
		return cutil.RespondMIME("", "application/json", "json", util.ToJSONBytes(rsp, true), rc)
	})
}
