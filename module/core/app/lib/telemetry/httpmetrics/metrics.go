package httpmetrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
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

func InjectHTTP(statusCode int, r *http.Request, span *telemetry.Span) {
	span.Attributes(
		&telemetry.Attribute{Key: "http.host", Value: r.Host},
		&telemetry.Attribute{Key: "http.method", Value: r.Method},
		&telemetry.Attribute{Key: "http.url", Value: r.URL.String()},
		&telemetry.Attribute{Key: "http.scheme", Value: r.URL.Scheme},
	)
	if s := r.Header.Get("User-Agent"); s != "" {
		span.Attribute("http.user_agent", s)
	}
	if s := r.Header.Get("Content-Length"); s != "" {
		span.Attribute("http.request_content_length", s)
	}
	span.SetHTTPStatus(statusCode)
}

func registerHTTPMetrics(logger util.Logger) {
	subsystem := "http"
	reqCnt = telemetry.MetricsCounter(subsystem, "requests_total", "The HTTP request counts processed.", logger, "key", "code", "method")
	reqDur = telemetry.MetricsHistogram(subsystem, "request_duration_seconds", "The HTTP request duration in seconds.", logger, "key", "code")
	reqSize = telemetry.MetricsSummary(subsystem, "request_size_bytes", "The HTTP request sizes in bytes.", logger)
	rspSize = telemetry.MetricsSummary(subsystem, "response_size_bytes", "The HTTP response sizes in bytes.", logger)
}
