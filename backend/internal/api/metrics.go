package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"monitoring-backend/internal/db"
)

type Metric struct {
	Time        time.Time `json:"time"`
	NodeName    string    `json:"node_name"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
}

// CollectMetricHandler handles POST requests from the monitoring agents to save metrics.
func CollectMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Use current time if the agent didn't provide one
	if metric.Time.IsZero() {
		metric.Time = time.Now()
	}

	query := `INSERT INTO metrics (time, node_name, cpu_usage, memory_usage, disk_usage) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Conn.Exec(context.Background(), query, metric.Time, metric.NodeName, metric.CPUUsage, metric.MemoryUsage, metric.DiskUsage)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetMetricsHandler handles GET requests to retrieve collected metrics.
func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch the last 100 metrics as a simple default implementation
	query := `SELECT time, node_name, cpu_usage, memory_usage, disk_usage FROM metrics ORDER BY time DESC LIMIT 100`
	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Time, &m.NodeName, &m.CPUUsage, &m.MemoryUsage, &m.DiskUsage); err != nil {
			continue // Skip problematic rows
		}
		metrics = append(metrics, m)
	}

	if metrics == nil {
		metrics = []Metric{} // Return empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
