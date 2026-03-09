package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
)

type MemoryStore struct {
	mu          sync.RWMutex
	users       []User
	defaultUser string
	user    string // active user
	metrics map[string]*metric.Metric
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metric.Metric),
	}
}

func (s *MemoryStore) AddUser(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if hasUser(s.users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserExists)
	}
	s.users = append(s.users, User{Name: name})
	if len(s.users) == 1 {
		s.defaultUser = name
	}
	return nil
}

func (s *MemoryStore) ListUsers() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users, nil
}

func (s *MemoryStore) DefaultUser() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.users) == 0 {
		return "", ErrNoUsers
	}
	if s.defaultUser == "" {
		return "", ErrNoDefaultUser
	}
	return s.defaultUser, nil
}

func (s *MemoryStore) SetDefaultUser(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !hasUser(s.users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserNotFound)
	}
	s.defaultUser = name
	return nil
}

func (s *MemoryStore) SetUser(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !hasUser(s.users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserNotFound)
	}
	s.user = name
	return nil
}

func (s *MemoryStore) key(metricName string) string {
	return s.user + "/" + metricName
}

func (s *MemoryStore) getOrCreate(name, unit string) *metric.Metric {
	k := s.key(name)
	m, ok := s.metrics[k]
	if !ok {
		m = &metric.Metric{Name: name, Unit: unit}
		s.metrics[k] = m
	}
	return m
}

func (s *MemoryStore) AddDataPoint(metricName string, unit string, dp metric.DataPoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.user == "" {
		return ErrNoUser
	}
	m := s.getOrCreate(metricName, unit)
	m.DataPoints = append(m.DataPoints, dp)
	return nil
}

func (s *MemoryStore) AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.user == "" {
		return ErrNoUser
	}
	m := s.getOrCreate(metricName, unit)
	m.AddItem(item, ts)
	return nil
}

func (s *MemoryStore) GetMetric(name string) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.user == "" {
		return nil, ErrNoUser
	}
	m, ok := s.metrics[s.key(name)]
	if !ok {
		return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	return m, nil
}

func (s *MemoryStore) GetMetricRange(name string, start, end time.Time) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.user == "" {
		return nil, ErrNoUser
	}
	m, ok := s.metrics[s.key(name)]
	if !ok {
		return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	return m.FilterRange(start, end), nil
}

func (s *MemoryStore) ListMetrics() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.user == "" {
		return nil, ErrNoUser
	}
	prefix := s.user + "/"
	var names []string
	for k := range s.metrics {
		if name, ok := strings.CutPrefix(k, prefix); ok {
			names = append(names, name)
		}
	}
	return names, nil
}
