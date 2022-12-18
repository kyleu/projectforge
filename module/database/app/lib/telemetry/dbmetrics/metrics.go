package dbmetrics

import (
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Metrics struct {
	Key     string
	stmtCnt *prometheus.CounterVec
	stmtDur *prometheus.HistogramVec
}

func NewMetrics(key string, db StatsGetter, logger util.Logger) (*Metrics, error) {
	m := &Metrics{Key: key}
	err := m.registerDatabaseMetrics("database", logger)
	if err != nil {
		return nil, err
	}
	err = prometheus.Register(newStatsCollector("database", db))
	if err != nil {
		return nil, errors.Wrap(err, "unable to register stats collector")
	}

	return m, nil
}

func (p *Metrics) registerDatabaseMetrics(subsystem string, logger util.Logger) error {
	stmtCntHelp := "The total number of SQL statements processed."
	p.stmtCnt = telemetry.MetricsCounter(subsystem, "statements_total", stmtCntHelp, logger, "database", "sql", "method")
	stmtDurHelp := "The SQL statement duration in seconds."
	p.stmtDur = telemetry.MetricsHistogram(subsystem, "statement_duration_seconds", stmtDurHelp, logger, "database", "sql", "method")

	return nil
}

func (p *Metrics) IncStmt(sql string, method string) {
	p.stmtCnt.WithLabelValues(p.Key, sql, method).Inc()
}

func (p *Metrics) CompleteStmt(q string, op string, started time.Time) {
	elapsed := float64(time.Since(started)) / float64(time.Second)
	p.stmtDur.WithLabelValues(p.Key, q, op).Observe(elapsed)
}

func (p *Metrics) Close() error {
	prometheus.Unregister(p.stmtCnt)
	prometheus.Unregister(p.stmtDur)
	return nil
}
