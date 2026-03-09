package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/arjungandhi/metrics/pkg/metric"
)

type LocalStore struct {
	dir string
}

func NewLocalStore(dir string) (*LocalStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating store dir: %w", err)
	}
	return &LocalStore{dir: dir}, nil
}

func (s *LocalStore) configDir() string {
	return filepath.Join(s.dir, "config")
}

func (s *LocalStore) configPath(key string) string {
	return filepath.Join(s.configDir(), key+".json")
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

func (s *LocalStore) DeleteMetric(name string) error {
	err := os.Remove(s.metricPath(name))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("metric %q: %w", name, ErrNotFound)
		}
		return err
	}
	return nil
}

func (s *LocalStore) SetConfig(key string, value json.RawMessage) error {
	if err := os.MkdirAll(s.configDir(), 0755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	return os.WriteFile(s.configPath(key), value, 0644)
}

func (s *LocalStore) GetConfig(key string) (json.RawMessage, error) {
	data, err := os.ReadFile(s.configPath(key))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("config %q: %w", key, ErrConfigNotFound)
		}
		return nil, err
	}
	return json.RawMessage(data), nil
}

func (s *LocalStore) DeleteConfig(key string) error {
	err := os.Remove(s.configPath(key))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("config %q: %w", key, ErrConfigNotFound)
		}
		return err
	}
	return nil
}

func (s *LocalStore) ListMetrics() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".json") {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return names, nil
}
