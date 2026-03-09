package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
)

// LocalStore persists metrics as JSON files in a directory.
// Each metric is stored as <dir>/<metric_name>.json.
type LocalStore struct {
	dir string
}

// NewLocalStore creates a LocalStore using the HEALTH_DIR env var.
// Creates the directory if it doesn't exist.
func NewLocalStore() (*LocalStore, error) {
	dir := os.Getenv("HEALTH_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting home directory: %w", err)
		}
		dir = filepath.Join(home, ".local", "share", "health")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating health dir: %w", err)
	}

	return &LocalStore{dir: dir}, nil
}

func (s *LocalStore) metricPath(name string) string {
	return filepath.Join(s.dir, name+".json")
}

func (s *LocalStore) loadMetric(name string) (*metric.Metric, error) {
	data, err := os.ReadFile(s.metricPath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("metric %q not found", name)
		}
		return nil, err
	}

	var m metric.Metric
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing metric %q: %w", name, err)
	}
	return &m, nil
}

func (s *LocalStore) saveMetric(m *metric.Metric) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.metricPath(m.Name), data, 0644)
}

func (s *LocalStore) AddDataPoint(metricName string, unit string, dp metric.DataPoint) error {
	m, err := s.loadMetric(metricName)
	if err != nil {
		// metric doesn't exist yet, create it
		m = &metric.Metric{Name: metricName, Unit: unit}
	}

	m.DataPoints = append(m.DataPoints, dp)
	return s.saveMetric(m)
}

func (s *LocalStore) AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error {
	m, err := s.loadMetric(metricName)
	if err != nil {
		m = &metric.Metric{Name: metricName, Unit: unit}
	}

	today := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	tomorrow := today.AddDate(0, 0, 1)

	// Find today's data point
	found := false
	for i := range m.DataPoints {
		dp := &m.DataPoints[i]
		if !dp.Time.Before(today) && dp.Time.Before(tomorrow) {
			dp.Items = append(dp.Items, item)
			dp.Value += item.Value
			found = true
			break
		}
	}

	if !found {
		m.DataPoints = append(m.DataPoints, metric.DataPoint{
			Time:  ts,
			Value: item.Value,
			Items: []metric.Item{item},
		})
	}

	return s.saveMetric(m)
}

func (s *LocalStore) GetMetric(name string) (*metric.Metric, error) {
	return s.loadMetric(name)
}

func (s *LocalStore) GetMetricRange(name string, start, end time.Time) (*metric.Metric, error) {
	m, err := s.loadMetric(name)
	if err != nil {
		return nil, err
	}

	filtered := &metric.Metric{
		Name: m.Name,
		Unit: m.Unit,
	}
	for _, dp := range m.DataPoints {
		if !dp.Time.Before(start) && !dp.Time.After(end) {
			filtered.DataPoints = append(filtered.DataPoints, dp)
		}
	}
	return filtered, nil
}

func (s *LocalStore) ListMetrics() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return names, nil
}
