package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	apiURL     string
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

func main() {
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://jsonplaceholder.typicode.com"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/data", handleData)
	http.HandleFunc("/health", handleHealth)

	log.Printf("backend-2 (e2e-test) listening on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleData(w http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.Get(apiURL + "/todos?_limit=5")
	if err != nil {
		http.Error(w, fmt.Sprintf("upstream error: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read error: %v", err), http.StatusInternalServerError)
		return
	}

	var todos []any
	json.Unmarshal(body, &todos)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"service": "backend-2",
		"source":  apiURL,
		"data":    todos,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
// e2e test: values.yaml flow 1771405412
