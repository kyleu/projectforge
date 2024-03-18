package clib

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

func WireRouter(r *router.Router, logger util.Logger) fasthttp.RequestHandler {
	AppRoutesList = r.List()
	p := httpmetrics.NewMetrics(util.AppKey, logger)
	if util.GetEnvBool("compression_enabled", false) {
		return fasthttp.CompressHandlerLevel(p.WrapHandler(r, true), fasthttp.CompressBestSpeed)
	}
	return p.WrapHandler(r, true)
}
