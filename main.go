package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const serverAddr = ":8080"

type RequestBody struct {
	KeyValue          string `json:"key_value"`
	SerialNumberValue string `json:"serial_number_value"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	go Task(reqBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "Task started"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/task", handleRequest)

	fmt.Printf("Starting server on %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
