package metric

import "time"

// Item represents a named component of a data point (e.g. a food in a meal)
type Item struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// DataPoint is a single timestamped measurement.
// For simple metrics (weight), just use Value.
// For composite metrics (calories), use Items to break it down (e.g. foods in a meal)
// and Value holds the total.
type DataPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
	Items []Item    `json:"items,omitempty"`
}

// Metric is a named health metric with its recorded data points
type Metric struct {
	Name       string      `json:"name"`
	Unit       string      `json:"unit"`
	DataPoints []DataPoint `json:"data_points"`
}
