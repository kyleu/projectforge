package cmenu

import (
	"context"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/graphql"
	"{{{ .Package }}}/app/lib/menu"
)

func graphQLMenu(_ context.Context, gqlSvc *graphql.Service) *menu.Item {
	l := gqlSvc.Keys()
	kids := make(menu.Items, 0, len(l))
	titles := gqlSvc.Titles()
	if len(l) > 1 {
		lo.ForEach(l, func(x string, _ int) {
			kids = append(kids, &menu.Item{Key: x, Title: titles[x], Description: "A GraphQL schema", Icon: "graph", Route: "/graphql/" + x})
		})
	}
	return &menu.Item{Key: "graphql", Title: "GraphQL", Description: "A graph-based API", Icon: "graph", Route: "/graphql", Children: kids}
}
