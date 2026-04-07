package models

import "time"

type Metric struct {
	Time        time.Time `json:"time"`
	NodeName    string    `json:"node_name" validate:"required"`
	CPUUsage    float64   `json:"cpu_usage" validate:"gte=0"`
	MemoryUsage float64   `json:"memory_usage" validate:"gte=0"`
	DiskUsage   float64   `json:"disk_usage" validate:"gte=0"`
}
