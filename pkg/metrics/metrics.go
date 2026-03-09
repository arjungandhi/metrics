package metrics

import (
	"os"
	"path/filepath"
	"time"

	"github.com/arjungandhi/metrics/pkg/config"
	"github.com/arjungandhi/metrics/pkg/metric"
	"github.com/arjungandhi/metrics/pkg/store"
)

// Client is the main entry point for interacting with metrics.
type Client struct {
	Store  store.Store
	Config *config.Config
}

// New creates a new Client. If cfg is nil, default config is used.
func New(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		var err error
		cfg, err = config.Default()
		if err != nil {
			return nil, err
		}
	}

	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return nil, err
	}

	s, err := store.NewSQLStore(filepath.Join(cfg.Dir, "metrics.db"))
	if err != nil {
		return nil, err
	}

	return &Client{Store: s, Config: cfg}, nil
}

// Close releases resources held by the client.
func (c *Client) Close() error {
	if closer, ok := c.Store.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// Add records a data point for the named metric.
func (c *Client) Add(name string, value float64, labels map[string]string) error {
	return c.Store.AddDataPoint(name, metric.DataPoint{
		Time:   time.Now(),
		Value:  value,
		Labels: labels,
	})
}

// AddAt records a data point for the named metric at a specific time.
func (c *Client) AddAt(name string, value float64, t time.Time, labels map[string]string) error {
	return c.Store.AddDataPoint(name, metric.DataPoint{
		Time:   t,
		Value:  value,
		Labels: labels,
	})
}

// Get returns all data points for a metric.
func (c *Client) Get(name string) (*metric.Metric, error) {
	return c.Store.GetMetric(name)
}

// List returns the names of all tracked metrics.
func (c *Client) List() ([]string, error) {
	return c.Store.ListMetrics()
}
