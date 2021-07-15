package controller

import (
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/assets"
)

const assetBase = "assets"

func Favicon(ctx *fasthttp.RequestCtx) {
	data, hash, contentType, err := assets.Asset(assetBase, "/favicon.ico")
	ZipResponse(ctx, data, hash, contentType, err)
}

func RobotsTxt(ctx *fasthttp.RequestCtx) {
	data, hash, contentType, err := assets.Asset(assetBase, "/robots.txt")
	ZipResponse(ctx, data, hash, contentType, err)
}

func Static(ctx *fasthttp.RequestCtx) {
	path, err := filepath.Abs(strings.TrimPrefix(string(ctx.Request.URI().Path()), "/assets"))
	if err == nil {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		data, hash, contentType, e := assets.Asset(assetBase, path)
		ZipResponse(ctx, data, hash, contentType, e)
	} else {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
	}
}

func ZipResponse(ctx *fasthttp.RequestCtx, data []byte, hash string, contentType string, err error) {
	if err == nil {
		ctx.Response.Header.Set("Content-Encoding", "gzip")
		ctx.Response.Header.SetContentType(contentType)
		// ctx.Response.Header.Add("Cache-Control", "public, max-age=31536000")
		ctx.Response.Header.Add("ETag", hash)
		if string(ctx.Request.Header.Peek("If-None-Match")) == hash {
			ctx.SetStatusCode(fasthttp.StatusNotModified)
		} else {
			ctx.SetStatusCode(fasthttp.StatusOK)
			_, _ = ctx.Write(data)
		}
	} else {
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
	}
}
