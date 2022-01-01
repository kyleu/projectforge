package dbmetrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type StatsGetter interface {
	Stats() sql.DBStats
}

type StatsCollector struct {
	sg                    StatsGetter
	maxOpenDesc           *prometheus.Desc
	openDesc              *prometheus.Desc
	inUseDesc             *prometheus.Desc
	idleDesc              *prometheus.Desc
	waitedForDesc         *prometheus.Desc
	blockedSecondsDesc    *prometheus.Desc
	closedMaxIdleDesc     *prometheus.Desc
	closedMaxLifetimeDesc *prometheus.Desc
	closedMaxIdleTimeDesc *prometheus.Desc
}

func newStatsCollector(subsystem string, dbName string, sg StatsGetter) *StatsCollector {
	namespace := ""
	labels := prometheus.Labels{"db_name": dbName}
	x := func(key string, help string) *prometheus.Desc {
		return prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, key), help, nil, labels)
	}
	return &StatsCollector{
		sg:                    sg,
		maxOpenDesc:           x("max_open", "Maximum number of open connections to the database."),
		openDesc:              x("open", "The number of established connections both in use and idle."),
		inUseDesc:             x("in_use", "The number of connections currently in use."),
		idleDesc:              x("idle", "The number of idle connections."),
		waitedForDesc:         x("waited_for", "The total number of connections waited for."),
		blockedSecondsDesc:    x("blocked_seconds", "The total time blocked waiting for a new connection."),
		closedMaxIdleDesc:     x("closed_max_idle", "The total number of connections closed due to SetMaxIdleConns."),
		closedMaxLifetimeDesc: x("closed_max_lifetime", "The total number of connections closed due to SetConnMaxLifetime."),
		closedMaxIdleTimeDesc: x("closed_max_idle_time", "The total number of connections closed due to SetConnMaxIdleTime."),
	}
}

// Describe implements the prometheus.Collector interface.
func (c *StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenDesc
	ch <- c.openDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitedForDesc
	ch <- c.blockedSecondsDesc
	ch <- c.closedMaxIdleDesc
	ch <- c.closedMaxLifetimeDesc
	ch <- c.closedMaxIdleTimeDesc
}

func (c *StatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.sg.Stats()
	ch <- prometheus.MustNewConstMetric(c.maxOpenDesc, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openDesc, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseDesc, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleDesc, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitedForDesc, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.blockedSecondsDesc, prometheus.CounterValue, stats.WaitDuration.Seconds())
	ch <- prometheus.MustNewConstMetric(c.closedMaxIdleDesc, prometheus.CounterValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.closedMaxLifetimeDesc, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
	ch <- prometheus.MustNewConstMetric(c.closedMaxIdleTimeDesc, prometheus.CounterValue, float64(stats.MaxIdleTimeClosed))
}
