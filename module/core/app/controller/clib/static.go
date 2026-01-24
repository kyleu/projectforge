package clib

import (
	"net/http"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/assets"
)

func Favicon(w http.ResponseWriter, _ *http.Request) {
	e, err := assets.Embed("favicon.ico")
	assetResponse(w, e, err)
}

func RobotsTxt(w http.ResponseWriter, _ *http.Request) {
	e, err := assets.Embed("robots.txt")
	assetResponse(w, e, err)
}

func Static(w http.ResponseWriter, r *http.Request) {
	p, _ := cutil.PathRichString(r, "path", false)
	if p.Contains("../") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("invalid path"))
	} else {
		e, err := assets.Embed(p.String())
		assetResponse(w, e, err)
	}
}

func assetResponse(w http.ResponseWriter, e *assets.Entry, err error) {
	if err == nil {
		w.Header().Set(cutil.HeaderContentType, e.Mime)
		w.Header().Set(cutil.HeaderCacheControl, "public, max-age=86400") // 24 hours
		w.WriteHeader(http.StatusOK)
		cutil.WriteCORS(w)
		_, _ = w.Write(e.Bytes)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
	}
}
