package provider

import (
	"github.com/arjungandhi/metrics/pkg/store"
)

// Provider is a source of metrics that can be synced.
type Provider interface {
	// Name returns the unique name of the provider (e.g. "strava").
	Name() string

	// Setup performs interactive configuration (OAuth, API keys, etc.)
	// and persists settings to the store.
	Setup(s store.Store) error

	// Sync fetches data from the source and writes metrics to the store.
	Sync(s store.Store) error
}
