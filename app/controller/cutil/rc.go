// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/app/util"
)

func RequestCtxToMap(rc *fasthttp.RequestCtx, data interface{}) map[string]interface{} {
	reqHeaders := make(map[string]string, rc.Request.Header.Len())
	rc.Request.Header.VisitAll(func(k, v []byte) {
		reqHeaders[string(k)] = string(v)
	})
	req := map[string]interface{}{
		"id":          rc.ID(),
		"url":         rc.URI().String(),
		"protocol":    string(rc.Request.URI().Scheme()),
		"host":        string(rc.Request.URI().Host()),
		"path":        string(rc.Request.URI().Path()),
		"queryString": string(rc.Request.URI().QueryString()),
		"headers":     reqHeaders,
		"bodySize":    len(rc.Request.Body()),
		"string":      rc.Request.String(),
	}
	rspHeaders := make(map[string]string, rc.Response.Header.Len())
	rc.Response.Header.VisitAll(func(k, v []byte) {
		rspHeaders[string(k)] = string(v)
	})
	rsp := map[string]interface{}{
		"code":     rc.Response.StatusCode(),
		"bodySize": len(rc.Response.Body()),
		"headers":  rspHeaders,
		"string":   rc.Response.String(),
	}
	return map[string]interface{}{"data": data, "request": req, "response": rsp}
}

func RequestCtxBool(rc *fasthttp.RequestCtx, key string) bool {
	return string(rc.URI().QueryArgs().Peek(key)) == util.BoolTrue
}
