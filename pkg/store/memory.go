package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
)

// MemoryStore is an in-memory implementation of Store
type MemoryStore struct {
	mu      sync.RWMutex
	metrics map[string]*metric.Metric
}

// NewMemoryStore creates a new empty MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metric.Metric),
	}
}

func (s *MemoryStore) AddDataPoint(metricName string, unit string, dp metric.DataPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.metrics[metricName]
	if !ok {
		m = &metric.Metric{Name: metricName, Unit: unit}
		s.metrics[metricName] = m
	}
	m.DataPoints = append(m.DataPoints, dp)
	return nil
}

func (s *MemoryStore) AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.metrics[metricName]
	if !ok {
		m = &metric.Metric{Name: metricName, Unit: unit}
		s.metrics[metricName] = m
	}

	today := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	tomorrow := today.AddDate(0, 0, 1)

	for i := range m.DataPoints {
		dp := &m.DataPoints[i]
		if !dp.Time.Before(today) && dp.Time.Before(tomorrow) {
			dp.Items = append(dp.Items, item)
			dp.Value += item.Value
			return nil
		}
	}

	m.DataPoints = append(m.DataPoints, metric.DataPoint{
		Time:  ts,
		Value: item.Value,
		Items: []metric.Item{item},
	})
	return nil
}

func (s *MemoryStore) GetMetric(name string) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %q not found", name)
	}
	return m, nil
}

func (s *MemoryStore) GetMetricRange(name string, start, end time.Time) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %q not found", name)
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

func (s *MemoryStore) ListMetrics() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	names := make([]string, 0, len(s.metrics))
	for name := range s.metrics {
		names = append(names, name)
	}
	return names, nil
}
