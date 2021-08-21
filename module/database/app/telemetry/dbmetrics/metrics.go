package dbmetrics

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	stmtCnt *prometheus.CounterVec
	stmtDur *prometheus.HistogramVec
}

func NewMetrics(dbName string, db StatsGetter) *Metrics {
	ss := strings.ReplaceAll(strings.ReplaceAll(dbName, "/", "_"), ".", "_")
	m := &Metrics{}
	m.registerDatabaseMetrics(ss)
	prometheus.MustRegister(newStatsCollector(ss, dbName, db))
	return m
}

func (p *Metrics) registerDatabaseMetrics(subsystem string) {
	cOpts := prometheus.CounterOpts{Subsystem: subsystem, Name: "statements_total", Help: "The total number of SQL statements processed."}
	p.stmtCnt = prometheus.NewCounterVec(cOpts, []string{"sql", "method"})

	hOpts := prometheus.HistogramOpts{Subsystem: subsystem, Name: "statement_duration_seconds", Help: "The SQL statement duration in seconds."}
	hOpts.Buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60}
	p.stmtDur = prometheus.NewHistogramVec(hOpts, []string{"sql", "method"})

	prometheus.MustRegister(p.stmtCnt, p.stmtDur)
}

func (p *Metrics) IncStmt(sql string, method string) {
	p.stmtCnt.WithLabelValues(sql, method).Inc()
}

func (p *Metrics) CompleteStmt(q string, op string, started time.Time, err error) {
	elapsed := float64(time.Since(started)) / float64(time.Second)
	p.stmtDur.WithLabelValues(q, op).Observe(elapsed)
}
