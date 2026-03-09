package web

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"time"

	"github.com/arjungandhi/health/pkg/metric"
	"github.com/arjungandhi/health/pkg/store"
)

//go:embed static/*
var staticFiles embed.FS

func Serve(addr string, s store.Store) error {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("GET /api/metrics", apiListMetrics(s))
	mux.HandleFunc("GET /api/metrics/{name}", apiGetMetric(s))
	mux.HandleFunc("POST /api/metrics/{name}/datapoints", apiAddDataPoint(s))
	mux.HandleFunc("POST /api/metrics/{name}/items", apiAddItem(s))

	return http.ListenAndServe(addr, mux)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func apiListMetrics(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := s.ListMetrics()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if names == nil {
			names = []string{}
		}
		writeJSON(w, names)
	}
}

func apiGetMetric(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		var m *metric.Metric
		var err error

		if startStr != "" && endStr != "" {
			start, e1 := time.ParseInLocation("2006-01-02", startStr, time.Local)
			end, e2 := time.ParseInLocation("2006-01-02", endStr, time.Local)
			if e1 != nil || e2 != nil {
				writeError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD)")
				return
			}
			end = end.Add(24*time.Hour - time.Nanosecond) // inclusive of whole day
			m, err = s.GetMetricRange(name, start, end)
		} else {
			m, err = s.GetMetric(name)
		}

		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeError(w, http.StatusNotFound, err.Error())
			} else {
				writeError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		writeJSON(w, m)
	}
}

type addDataPointReq struct {
	Value float64 `json:"value"`
	Date  string  `json:"date"` // YYYY-MM-DD, optional
	Unit  string  `json:"unit"`
}

func apiAddDataPoint(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")

		var req addDataPointReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		ts := time.Now()
		if req.Date != "" {
			var err error
			ts, err = time.ParseInLocation("2006-01-02", req.Date, time.Local)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid date: "+err.Error())
				return
			}
		}

		dp := metric.DataPoint{Time: ts, Value: req.Value}
		if err := s.AddDataPoint(name, req.Unit, dp); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, map[string]string{"status": "ok"})
	}
}

type addItemReq struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Date  string  `json:"date"`
	Unit  string  `json:"unit"`
}

func apiAddItem(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricName := r.PathValue("name")

		var req addItemReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "item name is required")
			return
		}

		ts := time.Now()
		if req.Date != "" {
			var err error
			ts, err = time.ParseInLocation("2006-01-02", req.Date, time.Local)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid date: "+err.Error())
				return
			}
		}

		item := metric.Item{Name: req.Name, Value: req.Value}
		if err := s.AddItemToDay(metricName, req.Unit, item, ts); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, map[string]string{"status": "ok"})
	}
}
