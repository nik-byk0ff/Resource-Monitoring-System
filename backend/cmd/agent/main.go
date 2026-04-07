package main

import (
	"log"
	"monitoring-backend/services"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/api"
	}

	agent := services.NewAgent("node-mock-01", apiURL)
	log.Printf("Starting mock agent node-mock-01 reporting to %s", apiURL)
	agent.Start(5) // Report every 5 seconds
}
