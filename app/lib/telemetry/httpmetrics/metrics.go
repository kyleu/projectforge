// Content managed by Project Forge, see [projectforge.md] for details.
package httpmetrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/lib/telemetry"
)

type Metrics struct {
	reqCnt  *prometheus.CounterVec
	reqDur  *prometheus.HistogramVec
	reqSize prometheus.Summary
	rspSize prometheus.Summary

	MetricsPath string
}

func NewMetrics(subsystem string) *Metrics {
	ss := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(subsystem, "-", "_"), "/", "_"), ".", "_")
	m := &Metrics{MetricsPath: defaultMetricPath}
	m.registerHTTPMetrics(ss)
	return m
}

func InjectHTTP(rc *fasthttp.RequestCtx, span *telemetry.Span) {
	span.Attributes(
		&telemetry.Attribute{Key: "http.host", Value: string(rc.Request.Host())},
		&telemetry.Attribute{Key: "http.method", Value: string(rc.Method())},
		&telemetry.Attribute{Key: "http.url", Value: string(rc.Request.RequestURI())},
		&telemetry.Attribute{Key: "http.scheme", Value: string(rc.Request.URI().Scheme())},
	)
	if b := rc.Request.Header.Peek("User-Agent"); len(b) > 0 {
		span.Attribute("http.user_agent", string(b))
	}
	if b := rc.Request.Header.Peek("Content-Length"); len(b) > 0 {
		span.Attribute("http.request_content_length", string(b))
	}
	span.SetHTTPStatus(rc.Response.StatusCode())
}

func (p *Metrics) registerHTTPMetrics(subsystem string) {
	cOpts := prometheus.CounterOpts{Subsystem: subsystem, Name: "requests_total", Help: "The HTTP request counts processed."}
	p.reqCnt = prometheus.NewCounterVec(cOpts, []string{"code", "method"})

	hOpts := prometheus.HistogramOpts{Subsystem: subsystem, Name: "request_duration_seconds", Help: "The HTTP request duration in seconds."}
	hOpts.Buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60}
	p.reqDur = prometheus.NewHistogramVec(hOpts, []string{"code"})

	reqOpts := prometheus.SummaryOpts{Subsystem: subsystem, Name: "request_size_bytes", Help: "The HTTP request sizes in bytes."}
	p.reqSize = prometheus.NewSummary(reqOpts)

	rspOpts := prometheus.SummaryOpts{Subsystem: subsystem, Name: "response_size_bytes", Help: "The HTTP response sizes in bytes."}
	p.rspSize = prometheus.NewSummary(rspOpts)

	prometheus.MustRegister(p.reqCnt, p.reqDur, p.reqSize, p.rspSize)
}
