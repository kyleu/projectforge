// Content managed by Project Forge, see [projectforge.md] for details.
package httpmetrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

var (
	reqCnt  *prometheus.CounterVec
	reqDur  *prometheus.HistogramVec
	reqSize prometheus.Summary
	rspSize prometheus.Summary

	metricsMu sync.Mutex
)

type Metrics struct {
	Key         string
	MetricsPath string
}

func NewMetrics(key string, logger util.Logger) *Metrics {
	m := &Metrics{Key: key, MetricsPath: defaultMetricPath}
	if reqCnt == nil {
		metricsMu.Lock()
		if reqCnt == nil {
			registerHTTPMetrics(logger)
		}
		metricsMu.Unlock()
	}
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

func registerHTTPMetrics(logger util.Logger) {
	subsystem := "http"
	reqCnt = telemetry.MetricsCounter(subsystem, "requests_total", "The HTTP request counts processed.", logger, "key", "code", "method")
	reqDur = telemetry.MetricsHistogram(subsystem, "request_duration_seconds", "The HTTP request duration in seconds.", logger, "key", "code")
	reqSize = telemetry.MetricsSummary(subsystem, "request_size_bytes", "The HTTP request sizes in bytes.", logger)
	rspSize = telemetry.MetricsSummary(subsystem, "response_size_bytes", "The HTTP response sizes in bytes.", logger)
}
