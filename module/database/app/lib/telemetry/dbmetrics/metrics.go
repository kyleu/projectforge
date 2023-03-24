package dbmetrics

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

var (
	stmtCnt *prometheus.CounterVec
	stmtDur *prometheus.HistogramVec
	stmtMu  sync.Mutex

	SkipDBMetrics = util.GetEnvBool("db_metrics_disabled", false)
)

type Metrics struct {
	Key string
}

func NewMetrics(key string, db StatsGetter, logger util.Logger) (*Metrics, error) {
	if stmtCnt == nil {
		stmtMu.Lock()
		if stmtCnt == nil {
			metricsKey := "database"
			err := registerDatabaseMetrics(metricsKey, logger)
			if err != nil {
				return nil, err
			}
			err = prometheus.Register(newStatsCollector(metricsKey, db))
			if err != nil {
				return nil, errors.Wrap(err, "unable to register stats collector")
			}
		}
		stmtMu.Unlock()
	}

	return &Metrics{Key: key}, nil
}

func (m *Metrics) IncStmt(sql string, method string) {
	if !SkipDBMetrics {
		stmtCnt.WithLabelValues(m.Key, sql, method).Inc()
	}
}

func (m *Metrics) CompleteStmt(q string, op string, started time.Time) {
	if !SkipDBMetrics {
		elapsed := float64(time.Since(started)) / float64(time.Second)
		stmtDur.WithLabelValues(m.Key, q, op).Observe(elapsed)
	}
}

func (m *Metrics) Close() error {
	return nil
}

func registerDatabaseMetrics(subsystem string, logger util.Logger) error {
	if !SkipDBMetrics {
		stmtCntHelp := "The total number of SQL statements processed."
		stmtCnt = telemetry.MetricsCounter(subsystem, "statements_total", stmtCntHelp, logger, "database", "sql", "method")
		stmtDurHelp := "The SQL statement duration in seconds."
		stmtDur = telemetry.MetricsHistogram(subsystem, "statement_duration_seconds", stmtDurHelp, logger, "database", "sql", "method")
	}
	return nil
}
