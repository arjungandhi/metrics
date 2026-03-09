package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
)

type MemoryStore struct {
	mu      sync.RWMutex
	metrics map[string]*metric.Metric
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metric.Metric),
	}
}

func (s *MemoryStore) AddDataPoint(metricName string, dp metric.DataPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.metrics[metricName]
	if !ok {
		m = &metric.Metric{Name: metricName}
		s.metrics[metricName] = m
	}

	m.DataPoints = append(m.DataPoints, dp)
	return nil
}

func (s *MemoryStore) GetMetric(name string) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	return m, nil
}

func (s *MemoryStore) GetMetricRange(name string, start, end time.Time) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	return m.FilterRange(start, end), nil
}

func (s *MemoryStore) ListMetrics() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var names []string
	for name := range s.metrics {
		names = append(names, name)
	}
	return names, nil
}
