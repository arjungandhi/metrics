package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arjungandhi/metrics/pkg/metric"
	_ "modernc.org/sqlite"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(dsn string) (*SQLStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS data_points (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			metric     TEXT    NOT NULL,
			time       TEXT    NOT NULL,
			value      REAL    NOT NULL,
			labels     TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_data_points_metric ON data_points(metric);
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	return &SQLStore{db: db}, nil
}

func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) AddDataPoint(metricName string, dp metric.DataPoint) error {
	var labelsJSON []byte
	if len(dp.Labels) > 0 {
		var err error
		labelsJSON, err = json.Marshal(dp.Labels)
		if err != nil {
			return fmt.Errorf("marshaling labels: %w", err)
		}
	}

	_, err := s.db.Exec(
		`INSERT INTO data_points (metric, time, value, labels) VALUES (?, ?, ?, ?)`,
		metricName,
		dp.Time.Format(time.RFC3339),
		dp.Value,
		labelsJSON,
	)
	return err
}

func (s *SQLStore) GetMetric(name string) (*metric.Metric, error) {
	rows, err := s.db.Query(
		`SELECT time, value, labels FROM data_points WHERE metric = ? ORDER BY time`,
		name,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := &metric.Metric{Name: name}
	for rows.Next() {
		var (
			timeStr   string
			value     float64
			labelsRaw []byte
		)
		if err := rows.Scan(&timeStr, &value, &labelsRaw); err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, fmt.Errorf("parsing time %q: %w", timeStr, err)
		}

		dp := metric.DataPoint{Time: t, Value: value}
		if len(labelsRaw) > 0 {
			if err := json.Unmarshal(labelsRaw, &dp.Labels); err != nil {
				return nil, fmt.Errorf("parsing labels: %w", err)
			}
		}
		m.DataPoints = append(m.DataPoints, dp)
	}

	if len(m.DataPoints) == 0 {
		return nil, fmt.Errorf("metric %q: %w", name, ErrNotFound)
	}
	return m, rows.Err()
}

func (s *SQLStore) ListMetrics() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT metric FROM data_points ORDER BY metric`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}
