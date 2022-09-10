package dbmetrics

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type Metrics struct {
	stmtCnt *prometheus.CounterVec
	stmtDur *prometheus.HistogramVec
}

func NewMetrics(key string, db StatsGetter, logger util.Logger) (*Metrics, error) {
	ss := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), "/", "_"), ".", "_")
	m := &Metrics{}
	err := m.registerDatabaseMetrics(ss, logger)
	if err != nil {
		return nil, err
	}
	err = prometheus.Register(newStatsCollector(ss, key, db))
	if err != nil {
		return nil, errors.Wrap(err, "unable to register stats collector")
	}

	return m, nil
}

func (p *Metrics) registerDatabaseMetrics(subsystem string, logger util.Logger) error {
	stmtCntHelp := "The total number of SQL statements processed."
	p.stmtCnt = telemetry.MetricsCounter(subsystem, "statements_total", stmtCntHelp, logger, "sql", "method")
	stmtDurHelp := "The SQL statement duration in seconds."
	p.stmtDur = telemetry.MetricsHistogram(subsystem, "statement_duration_seconds", stmtDurHelp, logger, "sql", "method")

	return nil
}

func (p *Metrics) IncStmt(sql string, method string) {
	p.stmtCnt.WithLabelValues(sql, method).Inc()
}

func (p *Metrics) CompleteStmt(q string, op string, started time.Time) {
	elapsed := float64(time.Since(started)) / float64(time.Second)
	p.stmtDur.WithLabelValues(q, op).Observe(elapsed)
}

func (p *Metrics) Close() error {
	prometheus.Unregister(p.stmtCnt)
	prometheus.Unregister(p.stmtDur)
	return nil
}
