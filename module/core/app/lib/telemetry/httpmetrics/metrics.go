package httpmetrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Metrics struct {
	reqCnt  *prometheus.CounterVec
	reqDur  *prometheus.HistogramVec
	reqSize prometheus.Summary
	rspSize prometheus.Summary

	MetricsPath string
}

func NewMetrics(subsystem string, logger util.Logger) *Metrics {
	ss := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(subsystem, "-", "_"), "/", "_"), ".", "_")
	m := &Metrics{MetricsPath: defaultMetricPath}
	m.registerHTTPMetrics(ss, logger)
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

func (p *Metrics) registerHTTPMetrics(subsystem string, logger util.Logger) {
	p.reqCnt = telemetry.MetricsCounter(subsystem, "requests_total", "The HTTP request counts processed.", logger, "code", "method")
	p.reqDur = telemetry.MetricsHistogram(subsystem, "request_duration_seconds", "The HTTP request duration in seconds.", logger, "code")
	p.reqSize = telemetry.MetricsSummary(subsystem, "request_size_bytes", "The HTTP request sizes in bytes.", logger)
	p.rspSize = telemetry.MetricsSummary(subsystem, "response_size_bytes", "The HTTP response sizes in bytes.", logger)
}
