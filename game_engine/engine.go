package game_engine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/glitchdawg/game-engine-with-user/api_server"
)

type GameEngine struct {
	winner           *api_server.UserResponse
	mu               sync.RWMutex
	totalResponses   int64
	correctResponses int64
	startTime        time.Time
	winnerFoundAt    *time.Time
	eventChan        chan GameEvent
	stopChan         chan bool
}

type GameEvent struct {
	Type     string
	Response api_server.UserResponse
	Time     time.Time
}

func NewGameEngine() *GameEngine {
	g := &GameEngine{
		startTime: time.Now(),
		eventChan: make(chan GameEvent, 1000),
		stopChan:  make(chan bool),
	}
	
	go g.processEvents()
	go g.printMetrics()
	
	return g
}

func (g *GameEngine) processEvents() {
	for {
		select {
		case event := <-g.eventChan:
			g.handleEvent(event)
		case <-g.stopChan:
			return
		}
	}
}

func (g *GameEngine) handleEvent(event GameEvent) {
	atomic.AddInt64(&g.totalResponses, 1)
	
	if event.Response.IsCorrect {
		atomic.AddInt64(&g.correctResponses, 1)
	}
	
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.winner == nil && event.Response.IsCorrect {
		g.winner = &event.Response
		now := time.Now()
		g.winnerFoundAt = &now
		
		timeTaken := now.Sub(g.startTime)
		
		fmt.Println("\n╔══════════════════════════════════════════╗")
		fmt.Println("║            🎉 WINNER FOUND! 🎉            ║")
		fmt.Println("╠══════════════════════════════════════════╣")
		fmt.Printf("║ Winner ID:     %-26d ║\n", event.Response.UserID)
		fmt.Printf("║ Answer:        %-26s ║\n", event.Response.Answer)
		fmt.Printf("║ Time to win:   %-26v ║\n", timeTaken)
		fmt.Printf("║ Total responses: %-24d ║\n", atomic.LoadInt64(&g.totalResponses))
		fmt.Printf("║ Correct answers: %-24d ║\n", atomic.LoadInt64(&g.correctResponses))
		fmt.Println("╚══════════════════════════════════════════╝\n")
	}
}

func (g *GameEngine) printMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			total := atomic.LoadInt64(&g.totalResponses)
			correct := atomic.LoadInt64(&g.correctResponses)
			
			if total > 0 {
				g.mu.RLock()
				hasWinner := g.winner != nil
				g.mu.RUnlock()
				
				if !hasWinner {
					percentage := float64(correct) / float64(total) * 100
					fmt.Printf("📊 Live Stats | Total: %d | Correct: %d (%.1f%%) | Duration: %v\n", 
						total, correct, percentage, time.Since(g.startTime).Round(time.Second))
				}
			}
		case <-g.stopChan:
			return
		}
	}
}

func (g *GameEngine) ProcessResponse(response api_server.UserResponse) bool {
	event := GameEvent{
		Type:     "response",
		Response: response,
		Time:     time.Now(),
	}
	
	select {
	case g.eventChan <- event:
	default:
		fmt.Println("Warning: Event channel full, processing synchronously")
		g.handleEvent(event)
	}
	
	g.mu.RLock()
	isWinner := g.winner != nil && g.winner.UserID == response.UserID
	g.mu.RUnlock()
	
	return isWinner
}

func (g *GameEngine) GetWinner() *api_server.UserResponse {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	if g.winner == nil {
		return nil
	}
	
	winnerCopy := *g.winner
	return &winnerCopy
}

func (g *GameEngine) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║           GAME ENGINE RESET              ║")
	fmt.Println("╠══════════════════════════════════════════╣")
	
	if g.winner != nil {
		fmt.Printf("║ Previous winner: User %-18d ║\n", g.winner.UserID)
	}
	
	total := atomic.LoadInt64(&g.totalResponses)
	correct := atomic.LoadInt64(&g.correctResponses)
	
	fmt.Printf("║ Total responses: %-23d ║\n", total)
	fmt.Printf("║ Correct responses: %-21d ║\n", correct)
	
	if total > 0 {
		percentage := float64(correct) / float64(total) * 100
		fmt.Printf("║ Success rate: %.1f%%                      ║\n", percentage)
	}
	
	fmt.Println("╚══════════════════════════════════════════╝\n")
	
	g.winner = nil
	atomic.StoreInt64(&g.totalResponses, 0)
	atomic.StoreInt64(&g.correctResponses, 0)
	g.startTime = time.Now()
	g.winnerFoundAt = nil
}

func (g *GameEngine) GetStats() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	total := atomic.LoadInt64(&g.totalResponses)
	correct := atomic.LoadInt64(&g.correctResponses)
	
	stats := map[string]interface{}{
		"total_responses":   total,
		"correct_responses": correct,
		"has_winner":        g.winner != nil,
		"game_duration":     time.Since(g.startTime).Seconds(),
	}
	
	if g.winner != nil {
		stats["winner_user_id"] = g.winner.UserID
		stats["winner_answer"] = g.winner.Answer
		if g.winnerFoundAt != nil {
			stats["time_to_win"] = g.winnerFoundAt.Sub(g.startTime).Seconds()
		}
	}
	
	if total > 0 {
		stats["correct_percentage"] = float64(correct) / float64(total) * 100
	}
	
	return stats
}

func (g *GameEngine) Shutdown() {
	close(g.stopChan)
}