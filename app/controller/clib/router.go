// Package clib - Content managed by Project Forge, see [projectforge.md] for details.
package clib

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/lib/telemetry/httpmetrics"
	"projectforge.dev/projectforge/app/util"
)

func WireRouter(r *router.Router, logger util.Logger) fasthttp.RequestHandler {
	AppRoutesList = r.List()
	p := httpmetrics.NewMetrics(util.AppKey, logger)
	if util.GetEnvBool("compression_enabled", false) {
		return fasthttp.CompressHandlerLevel(p.WrapHandler(r, true), fasthttp.CompressBestSpeed)
	}
	return p.WrapHandler(r, true)
}
