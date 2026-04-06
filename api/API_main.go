package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Transcript struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// Sample data - in real app, this would come from database
var transcripts = map[int]Transcript{
	1: {ID: 1, Content: "Transcript content for ticket 1...\nUser: Hello\nStaff: Hi there! How can I help you?\nUser: I have an issue with my account.\nStaff: Please provide more details."},
	2: {ID: 2, Content: "Transcript content for ticket 2...\nUser: Issue with login\nStaff: Please try resetting your password.\nUser: Done, but still not working.\nStaff: Let me check the logs."},
	3: {ID: 3, Content: "Transcript content for ticket 3...\nUser: Feature request\nStaff: What feature would you like?\nUser: Dark mode please.\nStaff: We'll consider it for the next update."},
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("../Dashboard/"))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/api/transcripts/", handleTranscripts)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTranscripts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/transcripts/")
	if path == "" {
		// Return all transcripts
		var allTranscripts []Transcript
		for _, t := range transcripts {
			allTranscripts = append(allTranscripts, t)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allTranscripts)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transcript, exists := transcripts[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Transcript not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transcript)
}
