package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/valyala/fasthttp"
)

type contextKey struct{}

var AuthorizationKey = &contextKey{}

func Auth(requestHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authorization := string(ctx.Request.Header.Peek("Authorization"))
		if authorization == "" {
			requestHandler(ctx)
		} else {
			// put raw authorization header value in context
			ctx.SetUserValue(AuthorizationKey, authorization)
			requestHandler(ctx)
		}
	}
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
func containsPath(ctx context.Context, path string) bool {
	for _, s := range GetPreloads(ctx) {
		if path == s {
			return true
		}
	}
	return false
}
