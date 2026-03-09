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

type localConfig struct {
	DefaultUser string `json:"default_user"`
	Users       []User `json:"users"`
}

// LocalStore persists metrics as JSON files under <dir>/users/<name>/<metric>.json.
type LocalStore struct {
	dir  string
	cfg  *localConfig
	user string // active user
}

func NewLocalStore(dir string) (*LocalStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating store dir: %w", err)
	}

	s := &LocalStore{dir: dir}
	if err := s.loadConfig(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *LocalStore) configPath() string {
	return filepath.Join(s.dir, "config.json")
}

func (s *LocalStore) loadConfig() error {
	data, err := os.ReadFile(s.configPath())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			s.cfg = &localConfig{}
			return nil
		}
		return fmt.Errorf("reading config: %w", err)
	}

	var cfg localConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	s.cfg = &cfg
	return nil
}

func (s *LocalStore) saveConfig() error {
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.configPath(), data, 0644)
}

func (s *LocalStore) AddUser(name string) error {
	if hasUser(s.cfg.Users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserExists)
	}
	s.cfg.Users = append(s.cfg.Users, User{Name: name})
	if len(s.cfg.Users) == 1 {
		s.cfg.DefaultUser = name
	}
	return s.saveConfig()
}

func (s *LocalStore) ListUsers() ([]User, error) {
	return s.cfg.Users, nil
}

func (s *LocalStore) DefaultUser() (string, error) {
	if len(s.cfg.Users) == 0 {
		return "", ErrNoUsers
	}
	if s.cfg.DefaultUser == "" {
		return "", ErrNoDefaultUser
	}
	return s.cfg.DefaultUser, nil
}

func (s *LocalStore) SetDefaultUser(name string) error {
	if !hasUser(s.cfg.Users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserNotFound)
	}
	s.cfg.DefaultUser = name
	return s.saveConfig()
}

func (s *LocalStore) SetUser(name string) error {
	if !hasUser(s.cfg.Users, name) {
		return fmt.Errorf("user %q: %w", name, ErrUserNotFound)
	}
	d := filepath.Join(s.dir, "users", name)
	if err := os.MkdirAll(d, 0755); err != nil {
		return err
	}
	s.user = name
	return nil
}

func (s *LocalStore) userDir() (string, error) {
	if s.user == "" {
		return "", ErrNoUser
	}
	return filepath.Join(s.dir, "users", s.user), nil
}

func (s *LocalStore) metricPath(name string) (string, error) {
	d, err := s.userDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, name+".json"), nil
}

func (s *LocalStore) loadMetric(name string) (*metric.Metric, error) {
	p, err := s.metricPath(name)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
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
	p, err := s.metricPath(m.Name)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func (s *LocalStore) AddDataPoint(metricName string, unit string, dp metric.DataPoint) error {
	m, err := s.loadMetric(metricName)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
		m = &metric.Metric{Name: metricName, Unit: unit}
	}

	m.DataPoints = append(m.DataPoints, dp)
	return s.saveMetric(m)
}

func (s *LocalStore) AddItemToDay(metricName string, unit string, item metric.Item, ts time.Time) error {
	m, err := s.loadMetric(metricName)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}
		m = &metric.Metric{Name: metricName, Unit: unit}
	}

	m.AddItem(item, ts)
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
	d, err := s.userDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(d)
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
