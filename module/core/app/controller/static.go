package controller

import (
	"strings"

	"github.com/valyala/fasthttp"

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
		rc.SetStatusCode(fasthttp.StatusOK)
		_, _ = rc.Write(data)
	} else {
		rc.Error(err.Error(), fasthttp.StatusNotFound)
	}
}
