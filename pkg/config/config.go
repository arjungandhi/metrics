package config

import (
	"os"
	"path/filepath"
)

// Config holds the configuration for a metrics instance.
type Config struct {
	// Dir is the data directory for storing metrics.
	// If empty, defaults to $METRICS_DIR or ~/.local/share/metrics.
	Dir string
}

// Default returns a Config with values resolved from the environment.
func Default() (*Config, error) {
	dir := os.Getenv("METRICS_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, ".local", "share", "metrics")
	}
	return &Config{Dir: dir}, nil
}
