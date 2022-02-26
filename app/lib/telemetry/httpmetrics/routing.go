// Content managed by Project Forge, see [projectforge.md] for details.
package httpmetrics

import (
	"strconv"
	"sync"
	"time"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	defaultMetricPath  = "/metrics"
	requestHandlerPool sync.Pool
)

func prometheusHandler() fasthttp.RequestHandler {
	return fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
}

func (p *Metrics) WrapHandler(r *router.Router) fasthttp.RequestHandler {
	r.GET(p.MetricsPath, prometheusHandler())

	return func(rc *fasthttp.RequestCtx) {
		if string(rc.Request.URI().Path()) == defaultMetricPath {
			r.Handler(rc)
			return
		}

		reqSize := make(chan int)
		frc := acquireRequestFromPool()
		rc.Request.CopyTo(frc)
		go computeApproximateRequestSize(frc, reqSize)

		start := time.Now()
		r.Handler(rc)

		status := strconv.Itoa(rc.Response.StatusCode())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		rspSize := float64(len(rc.Response.Body()))

		p.reqDur.WithLabelValues(status).Observe(elapsed)
		p.reqCnt.WithLabelValues(status, string(rc.Method())).Inc()
		p.reqSize.Observe(float64(<-reqSize))
		p.rspSize.Observe(rspSize)
	}
}

func computeApproximateRequestSize(rc *fasthttp.Request, out chan int) {
	s := 0
	if rc.URI() != nil {
		s += len(rc.URI().Path())
		s += len(rc.URI().Host())
	}
	s += len(rc.Header.Method())
	s += len("HTTP/1.1")
	rc.Header.VisitAll(func(key []byte, value []byte) {
		if string(key) != "Host" {
			s += len(key) + len(value)
		}
	})
	if rc.Header.ContentLength() != -1 {
		s += rc.Header.ContentLength()
	}
	out <- s
}

func acquireRequestFromPool() *fasthttp.Request {
	rp := requestHandlerPool.Get()
	if rp == nil {
		return &fasthttp.Request{}
	}
	frc, _ := rp.(*fasthttp.Request)
	return frc
}
