package clib

import (
	"strings"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/assets"
)

func Favicon(rc *fasthttp.RequestCtx) {
	e, err := assets.Embed("favicon.ico")
	assetResponse(rc, e, err)
}

func RobotsTxt(rc *fasthttp.RequestCtx) {
	e, err := assets.Embed("robots.txt")
	assetResponse(rc, e, err)
}

func Static(rc *fasthttp.RequestCtx) {
	p := strings.TrimPrefix(string(rc.Request.URI().Path()), "/assets")
	p = strings.TrimPrefix(p, "/")
	if strings.Contains(p, "../") {
		rc.Error("invalid path", fasthttp.StatusNotFound)
	} else {
		e, err := assets.Embed(p)
		assetResponse(rc, e, err)
	}
}

func assetResponse(rc *fasthttp.RequestCtx, e *assets.Entry, err error) {
	if err == nil {
		rc.Response.Header.SetContentType(e.Mime)
		rc.Response.Header.Set("Cache-Control", "public, max-age=3600")
		rc.SetStatusCode(fasthttp.StatusOK)
		cutil.WriteCORS(rc)
		_, _ = rc.Write(e.Bytes)
	} else {
		rc.Error(err.Error(), fasthttp.StatusNotFound)
	}
}
