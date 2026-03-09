package store

import (
	"errors"
	"slices"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
)

var (
	ErrNotFound      = errors.New("metric not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrNoUsers       = errors.New("no users configured")
	ErrNoUser        = errors.New("no active user set")
	ErrNoDefaultUser = errors.New("no default user set")
)

func hasUser(users []User, name string) bool {
	return slices.ContainsFunc(users, func(u User) bool { return u.Name == name })
}

type User struct {
	Name string `json:"name"`
}

type Store interface {
	// User management
	AddUser(name string) error
	ListUsers() ([]User, error)
	DefaultUser() (string, error)
	SetDefaultUser(name string) error
	SetUser(name string) error

	// Metric operations (require an active user)
	AddDataPoint(metricName string, unit string, dp metric.DataPoint) error
	AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error
	GetMetric(name string) (*metric.Metric, error)
	GetMetricRange(name string, start, end time.Time) (*metric.Metric, error)
	ListMetrics() ([]string, error)
}
