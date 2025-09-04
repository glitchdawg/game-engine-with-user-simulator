package mock_engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
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

type MockEngine struct {
	apiURL string
	wg     sync.WaitGroup
}

func NewMockEngine(apiURL string) *MockEngine {
	return &MockEngine{
		apiURL: apiURL,
	}
}

func (m *MockEngine) SimulateUsers(numUsers int) {
	fmt.Printf("Starting simulation for %d users...\n", numUsers)
	startTime := time.Now()

	for i := 1; i <= numUsers; i++ {
		m.wg.Add(1)
		go m.simulateUser(i)
	}

	m.wg.Wait()
	fmt.Printf("All %d users have sent their responses. Time taken: %v\n", numUsers, time.Since(startTime))
}

func (m *MockEngine) simulateUser(userID int) {
	defer m.wg.Done()

	isCorrect := rand.Float32() < 0.3
	
	delay := time.Duration(rand.Intn(991)+10) * time.Millisecond
	time.Sleep(delay)

	response := UserResponse{
		UserID:    userID,
		Answer:    generateAnswer(isCorrect),
		IsCorrect: isCorrect,
		Timestamp: time.Now().UnixNano(),
	}

	if err := m.sendResponse(response); err != nil {
		fmt.Printf("User %d failed to send response: %v\n", userID, err)
	}
}

func (m *MockEngine) sendResponse(response UserResponse) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	resp, err := http.Post(m.apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

func generateAnswer(isCorrect bool) string {
	correctAnswers := []string{"42", "correct", "true", "yes"}
	incorrectAnswers := []string{"41", "wrong", "false", "no", "maybe", "unknown"}

	if isCorrect {
		return correctAnswers[rand.Intn(len(correctAnswers))]
	}
	return incorrectAnswers[rand.Intn(len(incorrectAnswers))]
}