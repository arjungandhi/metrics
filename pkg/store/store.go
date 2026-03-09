package store

import (
	"time"

	"github.com/arjungandhi/health/pkg/metric"
)

// Store defines the interface for persisting and retrieving health metrics
type Store interface {
	// AddDataPoint adds a single data point to a metric.
	// Unit is set on first add and ignored on subsequent adds.
	AddDataPoint(metricName string, unit string, dp metric.DataPoint) error

	// AddItemToDay adds an item to the data point for the given day.
	// If no data point exists for that day, one is created.
	// The item's value is added to the data point's total.
	AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error

	// GetMetric returns a metric and all its data points
	GetMetric(name string) (*metric.Metric, error)

	// GetMetricRange returns data points for a metric within a time range
	GetMetricRange(name string, start, end time.Time) (*metric.Metric, error)

	// ListMetrics returns the names of all tracked metrics
	ListMetrics() ([]string, error)
}
