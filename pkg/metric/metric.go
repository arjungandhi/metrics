package metric

import "time"

// DataPoint is a single timestamped measurement with optional labels.
type DataPoint struct {
	Time   time.Time         `json:"time"`
	Value  float64           `json:"value"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Metric is a named time series.
type Metric struct {
	Name       string      `json:"name"`
	DataPoints []DataPoint `json:"data_points"`
}

// FilterRange returns a copy of the metric containing only data points within [start, end].
func (m *Metric) FilterRange(start, end time.Time) *Metric {
	filtered := &Metric{Name: m.Name}
	for _, dp := range m.DataPoints {
		if !dp.Time.Before(start) && !dp.Time.After(end) {
			filtered.DataPoints = append(filtered.DataPoints, dp)
		}
	}
	return filtered
}
