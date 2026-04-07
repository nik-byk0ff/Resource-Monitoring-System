package main

import (
	"log"
	"net/http"
	"os"

	"monitoring-backend/internal/api"
	"monitoring-backend/internal/db"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize Database
	db.InitDB()
	defer db.CloseDB()

	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/api/login", api.LoginHandler)

	// Protected Routes (Assuming agent uses same auth for simplicity)
	mux.HandleFunc("/api/metrics/collect", api.AuthMiddleware(api.CollectMetricHandler))
	
	// Admin protected routes
	mux.HandleFunc("/api/metrics", api.RoleMiddleware("admin", api.GetMetricsHandler))
    
    // User route (could limit what they see, for now just show metrics too for demo)
    mux.HandleFunc("/api/metrics/user", api.RoleMiddleware("user", api.GetMetricsHandler))

	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
