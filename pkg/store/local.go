package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
)

// LocalStore persists metrics as JSON files under <dir>/<metric>.json.
type LocalStore struct {
	dir string
}

func NewLocalStore(dir string) (*LocalStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating store dir: %w", err)
	}
	return &LocalStore{dir: dir}, nil
}

func (s *LocalStore) metricPath(name string) string {
	return filepath.Join(s.dir, name+".json")
}

func (s *LocalStore) loadMetric(name string) (*metric.Metric, error) {
	data, err := os.ReadFile(s.metricPath(name))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
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

func (s *LocalStore) AddDataPoint(metricName string, dp metric.DataPoint) error {
	m, err := s.loadMetric(metricName)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
		m = &metric.Metric{Name: metricName}
	}

	m.DataPoints = append(m.DataPoints, dp)
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
	return m.FilterRange(start, end), nil
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
