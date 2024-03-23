package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	connection_manager "github.com/kozr/stalk/communication_manager"
	rsakey "github.com/kozr/stalk/rsakey"
)

var rotationService *rsakey.KeyRotationService

func main() {
	rsakey.Init()
	rotationService = rsakey.GetRotationService()
	rotationService.SetRotationInterval(time.Hour * 24)
	rotationService.SetMaxKeyAge(time.Hour * 24 * 2)
	err := rotationService.StartKeyRotation()
	if err != nil {
		panic(err)
	}

	// Setup HTTP server
	http.HandleFunc("/public-key", publicKeyHandler)
	http.HandleFunc("/establish-connection", establishConnectionHandler)
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func publicKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch the latest public keys in PEM format
	pemStrings := rotationService.GetAllPublicKeys()

	if len(pemStrings) == 0 {
		http.Error(w, "No public keys available", http.StatusNotFound)
		return
	}

	// Prepare the data
	data := map[string]interface{}{
		"latest_key": pemStrings[len(pemStrings)-1], // Assuming the last key is the latest
		"past_keys":  pemStrings[:len(pemStrings)-1],
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal keys", http.StatusInternalServerError)
		return
	}

	// Correctly set the Content-Type for JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func establishConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Upgrade the connection to a WebSocket connection
	conn, err := connection_manager.GetConnectionManager().Upgrade(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	// Handle the WebSocket connection
	go handleConnection(conn)
}
