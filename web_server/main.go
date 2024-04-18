package web_server

import (
	"encoding/json"
	"net/http"

	"github.com/kozr/stalk/comm_manager"
	db "github.com/kozr/stalk/database"
	rsakey "github.com/kozr/stalk/rsakey"
)

// Add this function definition
func UpdateUserOnlineStatus(userId string, status bool) {
	// Implement your logic here
}

func PublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch the latest public keys in PEM format
	pemStrings := rsakey.GetRotationService().GetAllPublicKeys()

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

func EstablishConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// actually, find it in the path params
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Upgrade the connection to a WebSocket connection
	conn, err := comm_manager.GetConnectionManager().Upgrade(w, r, r.Header)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	conn.SetUserId(userId)
	db.UpdateUserOnlineStatus(userId, true)

	// Handle the WebSocket connection
	ch := make(chan string)
	go comm_manager.HandleIncoming(ch, conn)
	go comm_manager.HandleOutgoing(ch, conn)
}
