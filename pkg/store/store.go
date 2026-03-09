package store

import (
	"errors"

	"github.com/arjungandhi/metrics/pkg/metric"
)

var ErrNotFound = errors.New("metric not found")

type Store interface {
	AddDataPoint(metricName string, dp metric.DataPoint) error
	GetMetric(name string) (*metric.Metric, error)
	ListMetrics() ([]string, error)
}
