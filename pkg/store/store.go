package store

import (
	"errors"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
)

// ErrNotFound is returned when a metric does not exist.
var ErrNotFound = errors.New("metric not found")

// Store defines the interface for persisting and retrieving health metrics.
type Store interface {
	AddDataPoint(metricName string, unit string, dp metric.DataPoint) error
	AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error
	GetMetric(name string) (*metric.Metric, error)
	GetMetricRange(name string, start, end time.Time) (*metric.Metric, error)
	ListMetrics() ([]string, error)
}
