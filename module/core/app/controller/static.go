package controller

import (
	"path/filepath"
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
	path, err := filepath.Abs(strings.TrimPrefix(string(rc.Request.URI().Path()), "/assets"))
	if err == nil {
		path = strings.TrimPrefix(path, "/")
		data, contentType, e := assets.EmbedAsset(path)
		assetResponse(rc, data, contentType, e)
	} else {
		rc.Error(err.Error(), fasthttp.StatusBadRequest)
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
