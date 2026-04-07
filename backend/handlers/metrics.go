package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"monitoring-backend/db"
	"monitoring-backend/models"

	"github.com/jackc/pgx/v5"
)

func AddMetric(w http.ResponseWriter, r *http.Request) {
	var metric models.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if metric.Time.IsZero() {
		metric.Time = time.Now()
	}

	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO metrics (time, node_name, cpu_usage, memory_usage, disk_usage) VALUES ($1, $2, $3, $4, $5)",
		metric.Time, metric.NodeName, metric.CPUUsage, metric.MemoryUsage, metric.DiskUsage)
	if err != nil {
		http.Error(w, "failed to insert metric", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}

	rows, err := db.Pool.Query(context.Background(),
		"SELECT time, node_name, cpu_usage, memory_usage, disk_usage FROM metrics ORDER BY time DESC LIMIT $1", limit)
	if err != nil {
		http.Error(w, "failed to fetch metrics", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	metrics, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Metric])
	if err != nil {
		http.Error(w, "failed to parse metrics", http.StatusInternalServerError)
		return
	}
    
    // Ensure we return an empty array instead of null
    if metrics == nil {
        metrics = []models.Metric{}
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
