package store

import (
	"errors"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
)

var ErrNotFound = errors.New("metric not found")

// MatchLabels returns true if every key/value in the filter exists in the target.
func MatchLabels(target, filter map[string]string) bool {
	for k, v := range filter {
		if target[k] != v {
			return false
		}
	}
	return true
}

type Store interface {
	// AddDataPoint adds a data point to the named metric, creating it if needed.
	AddDataPoint(metricName string, dp metric.DataPoint) error

	// GetMetric returns the full metric by name.
	GetMetric(name string) (*metric.Metric, error)

	// GetMetricRange returns data points within [start, end].
	GetMetricRange(name string, start, end time.Time) (*metric.Metric, error)

	// ListMetrics returns all metric names.
	ListMetrics() ([]string, error)
}
