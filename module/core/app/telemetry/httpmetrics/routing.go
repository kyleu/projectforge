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

	return func(ctx *fasthttp.RequestCtx) {
		if string(ctx.Request.URI().Path()) == defaultMetricPath {
			r.Handler(ctx)
			return
		}

		reqSize := make(chan int)
		frc := acquireRequestFromPool()
		ctx.Request.CopyTo(frc)
		go computeApproximateRequestSize(frc, reqSize)

		start := time.Now()
		r.Handler(ctx)

		status := strconv.Itoa(ctx.Response.StatusCode())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		rspSize := float64(len(ctx.Response.Body()))

		p.reqDur.WithLabelValues(status).Observe(elapsed)
		p.reqCnt.WithLabelValues(status, string(ctx.Method())).Inc()
		p.reqSize.Observe(float64(<-reqSize))
		p.rspSize.Observe(rspSize)
	}
}

func computeApproximateRequestSize(ctx *fasthttp.Request, out chan int) {
	s := 0
	if ctx.URI() != nil {
		s += len(ctx.URI().Path())
		s += len(ctx.URI().Host())
	}
	s += len(ctx.Header.Method())
	s += len("HTTP/1.1")
	ctx.Header.VisitAll(func(key, value []byte) {
		if string(key) != "Host" {
			s += len(key) + len(value)
		}
	})
	if ctx.Header.ContentLength() != -1 {
		s += ctx.Header.ContentLength()
	}
	out <- s
}

func acquireRequestFromPool() *fasthttp.Request {
	rp := requestHandlerPool.Get()
	if rp == nil {
		return &fasthttp.Request{}
	}
	frc := rp.(*fasthttp.Request)
	return frc
}
