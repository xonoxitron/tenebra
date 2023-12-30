package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var urlCache []string
var cacheMutex sync.Mutex

// StartAPI initializes and starts the API server
func StartAPI() {
	// Load URLs into the cache on startup
	loadURLsIntoCache("tenebra.txt")

	r := mux.NewRouter()
	r.HandleFunc("/search", SearchHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "1991" // Default to port 1991 if not specified
	}

	logrus.Infof("API server listening on port %s", port)
	http.Handle("/", r)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Fatal("Error starting API server:", err)
	}
}

// SearchHandler handles search requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Query parameter 'query' is required", http.StatusBadRequest)
		return
	}

	// Search for URLs containing the query
	results := searchURLsInCache(query)

	// Respond with the results
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		logrus.Error("Error encoding JSON response:", err)
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

// loadURLsIntoCache loads URLs from the file into the cache
func loadURLsIntoCache(filePath string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urlCache = append(urlCache, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.Error("Error reading file:", err)
	}
}

// searchURLsInCache searches for URLs containing the query in the cache
func searchURLsInCache(query string) []string {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	var results []string

	for _, url := range urlCache {
		if strings.Contains(url, query) {
			results = append(results, url)
		}
	}

	return results
}
