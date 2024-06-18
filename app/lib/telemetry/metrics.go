package telemetry

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"projectforge.dev/projectforge/app/util"
)

var SkipControllerMetrics = util.GetEnvBool("controller_metrics_disabled", false)

const countSuffix = "Count"

type CounterAndHistogram struct {
	c *prometheus.CounterVec
	h *prometheus.HistogramVec
}

func NewCounterAndHistogram(subsystem string, name string, help string, logger util.Logger, labelNames ...string) *CounterAndHistogram {
	return &CounterAndHistogram{
		c: MetricsCounter(subsystem, name+countSuffix, help, logger, labelNames...),
		h: MetricsHistogram(subsystem, name, help, logger, labelNames...),
	}
}

func (x *CounterAndHistogram) Observe(startTime time.Time, labelValues ...string) {
	elapsed := float64(time.Since(startTime)) / float64(time.Second)
	x.h.WithLabelValues(labelValues...).Observe(elapsed)
	x.c.WithLabelValues(labelValues...).Inc()
}

type CounterAndGauge struct {
	c *prometheus.CounterVec
	g *prometheus.GaugeVec
}

func NewCounterAndGauge(subsystem string, name string, help string, logger util.Logger, labelNames ...string) *CounterAndGauge {
	return &CounterAndGauge{
		c: MetricsCounter(subsystem, name+countSuffix, help, logger, labelNames...),
		g: MetricsGauge(subsystem, name, help, logger, labelNames...),
	}
}

func (x *CounterAndGauge) Observe(labelValues ...string) {
	x.g.WithLabelValues(labelValues...).Inc()
	x.c.WithLabelValues(labelValues...).Inc()
}

func MetricsCounter(subsystem string, name string, help string, logger util.Logger, labelNames ...string) *prometheus.CounterVec {
	cOpts := prometheus.CounterOpts{Subsystem: subsystem, Name: name, Help: help}
	ret := prometheus.NewCounterVec(cOpts, labelNames)
	err := prometheus.Register(ret)
	if err != nil {
		logger.Warnf("error registering counter metric [%s:%s]: %+v", subsystem, name, err)
	}
	return ret
}

func MetricsGauge(subsystem string, name string, help string, logger util.Logger, labelNames ...string) *prometheus.GaugeVec {
	cOpts := prometheus.GaugeOpts{Subsystem: subsystem, Name: name, Help: help}
	ret := prometheus.NewGaugeVec(cOpts, labelNames)
	err := prometheus.Register(ret)
	if err != nil {
		logger.Warnf("error registering gauge metric [%s:%s]: %+v", subsystem, name, err)
	}
	return ret
}

func MetricsHistogram(subsystem string, name string, help string, logger util.Logger, labelNames ...string) *prometheus.HistogramVec {
	hOpts := prometheus.HistogramOpts{Subsystem: subsystem, Name: name, Help: help}
	hOpts.Buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60}
	ret := prometheus.NewHistogramVec(hOpts, labelNames)
	err := prometheus.Register(ret)
	if err != nil {
		logger.Warnf("error registering histogram metric [%s:%s]: %+v", subsystem, name, err)
	}
	return ret
}

func MetricsSummary(subsystem string, name string, help string, logger util.Logger) prometheus.Summary {
	reqOpts := prometheus.SummaryOpts{Subsystem: subsystem, Name: name, Help: help}
	ret := prometheus.NewSummary(reqOpts)
	err := prometheus.Register(ret)
	if err != nil {
		logger.Warnf("error registering summary metric [%s:%s]: %+v", subsystem, name, err)
	}
	return ret
}
