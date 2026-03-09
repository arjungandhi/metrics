package metrics

import (
	"testing"
	"time"

	"github.com/arjungandhi/metrics/pkg/config"
)

func newTestClient(t *testing.T, storeType config.StoreType) *Client {
	t.Helper()
	c, err := New(&config.Config{
		Dir:   t.TempDir(),
		Store: storeType,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

func TestAddAndGet(t *testing.T) {
	for _, st := range []config.StoreType{config.StoreSQL, config.StoreLocal, config.StoreMemory} {
		t.Run(string(st), func(t *testing.T) {
			c := newTestClient(t, st)

			if err := c.Add("weight", 185.5, nil); err != nil {
				t.Fatal(err)
			}

			m, err := c.Get("weight")
			if err != nil {
				t.Fatal(err)
			}
			if len(m.DataPoints) != 1 {
				t.Fatalf("got %d data points, want 1", len(m.DataPoints))
			}
			if m.DataPoints[0].Value != 185.5 {
				t.Errorf("value = %f, want 185.5", m.DataPoints[0].Value)
			}
		})
	}
}

func TestAddAtWithLabels(t *testing.T) {
	c := newTestClient(t, config.StoreMemory)

	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	labels := map[string]string{"unit": "lbs"}

	if err := c.AddAt("weight", 180.0, ts, labels); err != nil {
		t.Fatal(err)
	}

	m, err := c.Get("weight")
	if err != nil {
		t.Fatal(err)
	}
	dp := m.DataPoints[0]
	if !dp.Time.Equal(ts) {
		t.Errorf("time = %v, want %v", dp.Time, ts)
	}
	if dp.Labels["unit"] != "lbs" {
		t.Errorf("labels = %v, want unit=lbs", dp.Labels)
	}
}

func TestList(t *testing.T) {
	c := newTestClient(t, config.StoreMemory)

	c.Add("a", 1, nil)
	c.Add("b", 2, nil)

	names, err := c.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("got %d metrics, want 2", len(names))
	}
}

func TestGetNotFound(t *testing.T) {
	c := newTestClient(t, config.StoreMemory)

	_, err := c.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent metric")
	}
}
