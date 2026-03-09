package config

import (
	"testing"
)

func TestDefault(t *testing.T) {
	t.Setenv("METRICS_DIR", "/tmp/test-metrics")

	cfg, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Dir != "/tmp/test-metrics" {
		t.Errorf("Dir = %q, want /tmp/test-metrics", cfg.Dir)
	}
	if cfg.Store != StoreSQL {
		t.Errorf("Store = %q, want %q", cfg.Store, StoreSQL)
	}
}

func TestDefaultFallback(t *testing.T) {
	t.Setenv("METRICS_DIR", "")

	cfg, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Dir == "" {
		t.Error("Dir should not be empty")
	}
}
