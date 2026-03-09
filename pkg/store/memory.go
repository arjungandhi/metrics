package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
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

func (s *MemoryStore) getOrCreate(name, unit string) *metric.Metric {
	m, ok := s.metrics[name]
	if !ok {
		m = &metric.Metric{Name: name, Unit: unit}
		s.metrics[name] = m
	}
	return m
}

func (s *MemoryStore) AddDataPoint(metricName string, unit string, dp metric.DataPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.getOrCreate(metricName, unit)
	m.DataPoints = append(m.DataPoints, dp)
	return nil
}

func (s *MemoryStore) AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.getOrCreate(metricName, unit)
	m.AddItem(item, ts)
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

	names := make([]string, 0, len(s.metrics))
	for name := range s.metrics {
		names = append(names, name)
	}
	return names, nil
}
