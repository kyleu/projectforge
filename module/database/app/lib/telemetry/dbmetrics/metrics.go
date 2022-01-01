package dbmetrics

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	stmtCnt *prometheus.CounterVec
	stmtDur *prometheus.HistogramVec
}

func NewMetrics(key string, db StatsGetter) (*Metrics, error) {
	ss := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), "/", "_"), ".", "_")
	m := &Metrics{}
	err := m.registerDatabaseMetrics(ss)
	if err != nil {
		return nil, err
	}
	err = prometheus.Register(newStatsCollector(ss, key, db))
	if err != nil {
		return nil, errors.Wrap(err, "unable to register stats collector")
	}

	return m, nil
}

func (p *Metrics) registerDatabaseMetrics(subsystem string) error {
	cOpts := prometheus.CounterOpts{Subsystem: subsystem, Name: "statements_total", Help: "The total number of SQL statements processed."}

	p.stmtCnt = prometheus.NewCounterVec(cOpts, []string{"sql", "method"})
	err := prometheus.Register(p.stmtCnt)
	if err != nil {
		return errors.Wrap(err, "unable to register statement counter")
	}

	hOpts := prometheus.HistogramOpts{Subsystem: subsystem, Name: "statement_duration_seconds", Help: "The SQL statement duration in seconds."}
	hOpts.Buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60}
	p.stmtDur = prometheus.NewHistogramVec(hOpts, []string{"sql", "method"})
	err = prometheus.Register(p.stmtDur)
	if err != nil {
		return errors.Wrap(err, "unable to register statement duration histogram")
	}

	return nil
}

func (p *Metrics) IncStmt(sql string, method string) {
	p.stmtCnt.WithLabelValues(sql, method).Inc()
}

func (p *Metrics) CompleteStmt(q string, op string, started time.Time, err error) {
	elapsed := float64(time.Since(started)) / float64(time.Second)
	p.stmtDur.WithLabelValues(q, op).Observe(elapsed)
}

func (p *Metrics) Close() error {
	prometheus.Unregister(p.stmtCnt)
	prometheus.Unregister(p.stmtDur)
	return nil
}
