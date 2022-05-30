package controller

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/graphql"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vgraphql"
)

func GraphQLIndex(rc *fasthttp.RequestCtx) {
	act("graphql.index", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "GraphQL"
		keys := as.GraphQL.Keys()
		if len(keys) == 1 {
			return "/graphql/" + keys[0], nil
		}
		ps.Data = keys
		return render(rc, as, &vgraphql.List{Keys: keys}, ps, "graphql")
	})
}

func GraphQLDetail(rc *fasthttp.RequestCtx) {
	act("graphql.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		titles := as.GraphQL.Titles()
		if len(titles) > 1 {
			ps.Title = fmt.Sprintf("[%s] GraphQL", titles[key])
		} else {
			ps.Title = "GraphQL"
		}
		ps.Data = key
		bc := []string{"graphql"}
		if len(titles) > 1 {
			bc = append(bc, key)
		}
		return render(rc, as, &vgraphql.Detail{Key: key}, ps, bc...)
	})
}

func GraphQLRun(rc *fasthttp.RequestCtx) {
	act("graphql.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		sch := as.GraphQL.Schema(key)
		if sch == nil {
			return "", errors.Errorf("no schema found at [%s]", key)
		}

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse form")
		}
		q := frm.GetStringOpt("query")
		op := frm.GetStringOpt("operationName")
		v := frm.GetStringOpt("variables")
		vars := util.ValueMap{}
		if v != "" {
			_ = util.FromJSON([]byte(v), &vars)
		}
		rsp, err := as.GraphQL.Exec(ps.Context, sch, q, op, vars)
		if err != nil {
			return "", err
		}
		return cutil.RespondMIME("", "application/json", "json", util.ToJSONBytes(rsp, true), rc)
	})
}

func graphQLMenu(gqlSvc *graphql.Service, ctx context.Context) *menu.Item {
	l := gqlSvc.Keys()
	kids := make(menu.Items, 0, len(l))
	titles := gqlSvc.Titles()
	if len(l) > 1 {
		for _, x := range l {
			kids = append(kids, &menu.Item{Key: x, Title: titles[x], Description: "A GraphQL schema", Icon: "graph", Route: "/graphql/" + x})
		}
	}
	return &menu.Item{Key: "graphql", Title: "GraphQL", Description: "A graph-based API", Icon: "graph", Route: "/graphql", Children: kids}
}
