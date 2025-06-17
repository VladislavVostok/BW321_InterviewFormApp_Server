package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Модель данных (Candidate)
type Candidate struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullName"`
	Age             int    `json:"age"`
	DesiredSalary   int    `json:"desiredSalary"`
	CorrectAnswers  int    `json:"correctAnswers"`
	TotalScore      int    `json:"totalScore"`
	HasExperience   bool   `json:"hasExperience"`
	HasTeamSkills   bool   `json:"hasTeamSkills"`
	ReadyForTrips   bool   `json:"readyForTrips"`
	PassedInterview bool   `json:"passedInterview"`
}

// Хранилище данных (в памяти)
var (
	candidates = make(map[int]Candidate)
	nextID     = 1
	mu         sync.Mutex
)

// Извлекает ID из URL (например, "/candidates/42" → 42)
func extractIDFromURL(r *http.Request) (int, error) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid URL format")
	}
	idStr := parts[2]
	return strconv.Atoi(idStr)
}

// GET /candidates
func getCandidates(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

// POST /candidates
func addCandidate(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var candidate Candidate
	if err := json.NewDecoder(r.Body).Decode(&candidate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Отработал чётко!!!!!!!!!!!!")
	fmt.Println(nextID)
	candidate.ID = nextID
	nextID++
	candidates[candidate.ID] = candidate

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(candidate)
}

// GET /candidates/{id}
func getCandidate(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id, err := extractIDFromURL(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	candidate, exists := candidates[id]
	if !exists {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidate)
}

// DELETE /candidates/{id}
func deleteCandidate(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id, err := extractIDFromURL(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if _, exists := candidates[id]; !exists {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	delete(candidates, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/candidates", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCandidates(w, r)
		case http.MethodPost:
			addCandidate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/candidates/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCandidate(w, r)
		case http.MethodDelete:
			deleteCandidate(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := "9999"
	fmt.Printf("REST-сервер запущен на http://45.144.221.6:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
