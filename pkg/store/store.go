package store

import (
	"encoding/json"
	"errors"

	"github.com/arjungandhi/metrics/pkg/metric"
)

var ErrNotFound = errors.New("metric not found")
var ErrConfigNotFound = errors.New("config not found")

type Store interface {
	AddDataPoint(metricName string, dp metric.DataPoint) error
	GetMetric(name string) (*metric.Metric, error)
	ListMetrics() ([]string, error)

	// Config key-value storage for provider settings.
	// Values are stored as JSON.
	SetConfig(key string, value json.RawMessage) error
	GetConfig(key string) (json.RawMessage, error)
	DeleteConfig(key string) error
}
