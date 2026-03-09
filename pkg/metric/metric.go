package metric

import "time"

// Item is a named component of a data point (e.g. a food in a meal).
type Item struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// DataPoint is a single timestamped measurement. Items optionally break down
// the total into named components.
type DataPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
	Items []Item    `json:"items,omitempty"`
}

type Metric struct {
	Name       string      `json:"name"`
	Unit       string      `json:"unit"`
	DataPoints []DataPoint `json:"data_points"`
}

// AddItem accumulates an item into the day's data point, creating one if needed.
func (m *Metric) AddItem(item Item, ts time.Time) {
	today := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	tomorrow := today.AddDate(0, 0, 1)

	for i := range m.DataPoints {
		dp := &m.DataPoints[i]
		if !dp.Time.Before(today) && dp.Time.Before(tomorrow) {
			dp.Items = append(dp.Items, item)
			dp.Value += item.Value
			return
		}
	}

	m.DataPoints = append(m.DataPoints, DataPoint{
		Time:  ts,
		Value: item.Value,
		Items: []Item{item},
	})
}

// FilterRange returns a copy of the metric containing only data points within [start, end].
func (m *Metric) FilterRange(start, end time.Time) *Metric {
	filtered := &Metric{
		Name: m.Name,
		Unit: m.Unit,
	}
	for _, dp := range m.DataPoints {
		if !dp.Time.Before(start) && !dp.Time.After(end) {
			filtered.DataPoints = append(filtered.DataPoints, dp)
		}
	}
	return filtered
}
