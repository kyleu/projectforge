// Package httpmetrics - Content managed by Project Forge, see [projectforge.md] for details.
package httpmetrics

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"projectforge.dev/projectforge/app/util"
)

var defaultMetricPath = "/metrics"

func (p *Metrics) WrapHandler(router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBytes := make(chan int)
		go computeApproximateRequestSize(r, reqBytes)

		start := util.TimeCurrent()
		router.ServeHTTP(w, r)
		elapsed := float64(time.Since(start)) / float64(time.Second)
		// status := strconv.Itoa(rc.Response.StatusCode())
		status := "200"
		// rspBytes := float64(len(rc.Response.Body()))
		rspBytes := 100.0

		reqDur.WithLabelValues(p.Key, status).Observe(elapsed)
		reqCnt.WithLabelValues(p.Key, status, r.Method).Inc()
		reqSize.Observe(float64(<-reqBytes))
		rspSize.Observe(rspBytes)
	}
}

func computeApproximateRequestSize(r *http.Request, out chan int) {
	s := 0
	if r.URL != nil {
		s += len(r.URL.Path)
		s += len(r.URL.Host)
	}
	s += len(r.Method)
	s += len("HTTP/1.1")
	for k, v := range r.Header {
		if k != "Host" {
			s += len(k) + len(v)
		}
	}
	out <- s
}
