package clib

import (
	"net/http"
	"strings"

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
	p, _ := cutil.PathString(r, "path", false)
	if strings.Contains(p, "../") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("invalid path"))
	} else {
		e, err := assets.Embed(p)
		assetResponse(w, e, err)
	}
}

func assetResponse(w http.ResponseWriter, e *assets.Entry, err error) {
	if err == nil {
		w.Header().Set(cutil.HeaderContentType, e.Mime)
		w.Header().Set(cutil.HeaderCacheControl, "public, max-age=3600")
		w.WriteHeader(http.StatusOK)
		cutil.WriteCORS(w)
		_, _ = w.Write(e.Bytes)
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
	}
}
