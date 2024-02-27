package clib

import (
	"strings"

	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/assets"
)

func Favicon(rc *fasthttp.RequestCtx) {
	data, contentType, err := assets.EmbedAsset("favicon.ico")
	assetResponse(rc, data, contentType, err)
}

func RobotsTxt(rc *fasthttp.RequestCtx) {
	data, contentType, err := assets.EmbedAsset("robots.txt")
	assetResponse(rc, data, contentType, err)
}

func Static(rc *fasthttp.RequestCtx) {
	p := strings.TrimPrefix(string(rc.Request.URI().Path()), "/assets")
	p = strings.TrimPrefix(p, "/")
	if strings.Contains(p, "../") {
		rc.Error("invalid path", fasthttp.StatusNotFound)
	} else {
		data, contentType, e := assets.EmbedAsset(p)
		assetResponse(rc, data, contentType, e)
	}
}

func assetResponse(rc *fasthttp.RequestCtx, data []byte, contentType string, err error) {
	if err == nil {
		rc.Response.Header.SetContentType(contentType)
		if !util.DEBUG {
			rc.Response.Header.Set("Cache-Control", "public, max-age=300")
		}
		rc.SetStatusCode(fasthttp.StatusOK)
		cutil.WriteCORS(rc)
		_, _ = rc.Write(data)
	} else {
		rc.Error(err.Error(), fasthttp.StatusNotFound)
	}
}
