package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/kozr/stalk/database"
	redis_client "github.com/kozr/stalk/redis"
	rsakey "github.com/kozr/stalk/rsakey"
	web_server "github.com/kozr/stalk/web_server"
)

func main() {
	if err := initServices(); err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}
	defer cleanup()

	// Setup HTTP server
	http.HandleFunc("/public-key", web_server.PublicKeyHandler)
	http.HandleFunc("/establish-connection", web_server.EstablishConnectionHandler)
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServices() error {
	var err error
	if err = rsakey.Init(); err != nil {
		return fmt.Errorf("RSA Key initialization failed: %w", err)
	}
	if err = redis_client.Init(); err != nil {
		return fmt.Errorf("Redis client initialization failed: %w", err)
	}
	if err = database.Init(); err != nil {
		return fmt.Errorf("Database initialization failed: %w", err)
	}
	rotationService := rsakey.GetRotationService()
	rotationService.SetRotationInterval(time.Hour * 24)
	rotationService.SetMaxKeyAge(time.Hour * 24 * 2)
	if err = rotationService.StartKeyRotation(); err != nil {
		return fmt.Errorf("key rotation service failed to start: %w", err)
	}
	return nil
}

func cleanup() {
	database.DB.Close()
}
