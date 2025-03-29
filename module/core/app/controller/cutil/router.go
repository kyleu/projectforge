package cutil

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/CAFxX/httpcompression"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

var AppRoutesList = map[string][]string{}

func WireRouter(r *mux.Router, notFound http.HandlerFunc, logger util.Logger) (http.Handler, error) {
	p := httpmetrics.NewMetrics(util.AppKey, logger)
	r.Handle(p.MetricsPath, promhttp.Handler()).Methods(http.MethodGet)

	r.PathPrefix("/").HandlerFunc(notFound)

	var ret http.Handler = p.WrapHandler(r)
	includeCompression := util.GetEnvBool("compression_enabled", false)
	if includeCompression {
		compressedTypes := []string{
			"application/gzip", "application/octet-stream", "application/zip",
			"audio/aac", "audio/mpeg", "audio/ogg",
			"image/gif", "image/jpeg", "image/png", "image/webp",
			"video/mpeg", "video/mp4", "video/webm",
		}
		compress, err := httpcompression.DefaultAdapter(httpcompression.ContentTypes(compressedTypes, true))
		if err != nil {
			return nil, err
		}
		ret = compress(ret)
	}
	return ret, nil
}

func AddRoute(method string, path string) {
	curr := AppRoutesList[method]
	if !slices.Contains(curr, path) {
		AppRoutesList[method] = util.ArraySorted(append(curr, path))
	}
}

func URLAddQuery(u *url.URL, k string, v string) {
	q := u.Query()
	q.Set(k, v)
	u.RawQuery = q.Encode()
}
