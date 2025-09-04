package api_server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type UserResponse struct {
	UserID    int    `json:"user_id"`
	Answer    string `json:"answer"`
	IsCorrect bool   `json:"is_correct"`
	Timestamp int64  `json:"timestamp"`
}

type APIServer struct {
	port          string
	gameEngine    GameEngineInterface
	mu            sync.RWMutex
	totalReceived int
	startTime     time.Time
}

type GameEngineInterface interface {
	ProcessResponse(response UserResponse) bool
	GetWinner() *UserResponse
	Reset()
}

func NewAPIServer(port string, gameEngine GameEngineInterface) *APIServer {
	return &APIServer{
		port:       port,
		gameEngine: gameEngine,
		startTime:  time.Now(),
	}
}

func (s *APIServer) Start() error {
	http.HandleFunc("/submit", s.handleSubmit)

	log.Printf("API Server starting on port %s (endpoint: /submit)", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}

func (s *APIServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var response UserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.totalReceived++
	count := s.totalReceived
	s.mu.Unlock()

	isWinner := s.gameEngine.ProcessResponse(response)

	result := map[string]interface{}{
		"received":  true,
		"user_id":   response.UserID,
		"is_winner": isWinner,
		"response_count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

	if count%100 == 0 {
		log.Printf("Processed %d responses", count)
	}
}


func (s *APIServer) GetTotalResponses() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalReceived
}