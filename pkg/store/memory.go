package store

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/arjungandhi/metrics/pkg/metric"
)

type MemoryStore struct {
	mu      sync.RWMutex
	metrics map[string]*metric.Metric
	config  map[string]json.RawMessage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metric.Metric),
		config:  make(map[string]json.RawMessage),
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

	cp := &metric.Metric{
		Name:       m.Name,
		DataPoints: make([]metric.DataPoint, len(m.DataPoints)),
	}
	copy(cp.DataPoints, m.DataPoints)
	return cp, nil
}

func (s *MemoryStore) DeleteMetric(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.metrics[name]; !ok {
		return fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	delete(s.metrics, name)
	return nil
}

func (s *MemoryStore) SetConfig(key string, value json.RawMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config[key] = append(json.RawMessage(nil), value...)
	return nil
}

func (s *MemoryStore) GetConfig(key string) (json.RawMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.config[key]
	if !ok {
		return nil, fmt.Errorf("config %q: %w", key, ErrConfigNotFound)
	}
	cp := append(json.RawMessage(nil), v...)
	return cp, nil
}

func (s *MemoryStore) DeleteConfig(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.config[key]; !ok {
		return fmt.Errorf("config %q: %w", key, ErrConfigNotFound)
	}
	delete(s.config, key)
	return nil
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
