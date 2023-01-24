package graph

import "github.com/valyala/fasthttp"

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
