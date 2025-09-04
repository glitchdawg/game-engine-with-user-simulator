package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/glitchdawg/game-engine-with-user/api_server"
	"github.com/glitchdawg/game-engine-with-user/game_engine"
	"github.com/glitchdawg/game-engine-with-user/mock_engine"
)

func main() {
	var mode string
	var port string
	var numUsers int
	var apiURL string

	flag.StringVar(&mode, "mode", "server", "Mode: server, mock, or full")
	flag.StringVar(&port, "port", "8080", "API server port")
	flag.IntVar(&numUsers, "users", 1000, "Number of mock users")
	flag.StringVar(&apiURL, "api", "http://localhost:8080/submit", "API URL for mock engine")
	flag.Parse()

	switch mode {
	case "server":
		runInteractiveServer(port)
	case "mock":
		runMockEngine(numUsers, apiURL)
	case "full":
		runFullSimulation(port, numUsers)
	default:
		fmt.Println("Invalid mode. Use: server, mock, or full")
		os.Exit(1)
	}
}

func runInteractiveServer(port string) {
	clearScreen()
	printBanner("GAME SERVER")
	
	engine := game_engine.NewGameEngine()
	server := api_server.NewAPIServer(port, engine)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	time.Sleep(1 * time.Second)
	
	fmt.Printf("✅ Server running on port %s\n", port)
	fmt.Printf("📍 Endpoint: POST http://localhost:%s/submit\n", port)
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║         AVAILABLE COMMANDS         ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Println("║  stats  - Show current statistics  ║")
	fmt.Println("║  reset  - Reset the game engine    ║")
	fmt.Println("║  clear  - Clear the screen         ║")
	fmt.Println("║  exit   - Shutdown server          ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println("\n⏳ Waiting for user responses...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	go func() {
		<-sigChan
		handleShutdown(engine)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}
		
		command := strings.TrimSpace(strings.ToLower(scanner.Text()))
		
		switch command {
		case "stats":
			showStats(engine)
		case "reset":
			engine.Reset()
		case "clear":
			clearScreen()
			printBanner("GAME SERVER")
		case "exit", "quit":
			handleShutdown(engine)
		case "":
			continue
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}

func runMockEngine(numUsers int, apiURL string) {
	clearScreen()
	printBanner("MOCK USER ENGINE")
	
	fmt.Printf("📊 Simulating %d users\n", numUsers)
	fmt.Printf("🎯 Target API: %s\n", apiURL)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	engine := mock_engine.NewMockEngine(apiURL)
	start := time.Now()
	
	fmt.Println("\n⚡ Starting simulation...")
	engine.SimulateUsers(numUsers)
	
	duration := time.Since(start)
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║       SIMULATION COMPLETE          ║")
	fmt.Println("╠════════════════════════════════════╣")
	fmt.Printf("║ Total Users: %-22d ║\n", numUsers)
	fmt.Printf("║ Duration: %-25v ║\n", duration.Round(time.Millisecond))
	fmt.Printf("║ Avg time/user: %-20v ║\n", duration/time.Duration(numUsers))
	fmt.Println("╚════════════════════════════════════╝")
}

func runFullSimulation(port string, numUsers int) {
	clearScreen()
	printBanner("FULL SIMULATION")
	
	fmt.Printf("🖥️  Server Port: %s\n", port)
	fmt.Printf("👥 Mock Users: %d\n", numUsers)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	engine := game_engine.NewGameEngine()
	server := api_server.NewAPIServer(port, engine)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("\n✅ Server is ready")
	fmt.Println("⚡ Starting mock users...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mockEngine := mock_engine.NewMockEngine("http://localhost:" + port + "/submit")
	
	start := time.Now()
	mockEngine.SimulateUsers(numUsers)
	
	time.Sleep(2 * time.Second)
	
	displayFinalResults(engine, start)
}

func showStats(engine *game_engine.GameEngine) {
	stats := engine.GetStats()
	
	fmt.Println("\n╔════════════════════════════════════╗")
	fmt.Println("║         CURRENT STATISTICS        ║")
	fmt.Println("╠════════════════════════════════════╣")
	
	total := stats["total_responses"].(int64)
	correct := stats["correct_responses"].(int64)
	duration := stats["game_duration"].(float64)
	
	fmt.Printf("║ Total Responses: %-18d ║\n", total)
	fmt.Printf("║ Correct Responses: %-16d ║\n", correct)
	
	if total > 0 {
		percentage := stats["correct_percentage"].(float64)
		fmt.Printf("║ Success Rate: %.1f%%                ║\n", percentage)
	}
	
	fmt.Printf("║ Duration: %.1fs                     ║\n", duration)
	
	if stats["has_winner"].(bool) {
		fmt.Println("╠════════════════════════════════════╣")
		fmt.Printf("║ 🏆 Winner: User %-18d ║\n", stats["winner_user_id"])
		fmt.Printf("║ Answer: %-27s ║\n", stats["winner_answer"])
		timeToWin := stats["time_to_win"].(float64)
		fmt.Printf("║ Time to win: %.3fs                ║\n", timeToWin)
	} else {
		fmt.Println("╠════════════════════════════════════╣")
		fmt.Println("║ ⏳ No winner yet                   ║")
	}
	
	fmt.Println("╚════════════════════════════════════╝")
}

func displayFinalResults(engine *game_engine.GameEngine, startTime time.Time) {
	stats := engine.GetStats()
	
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║          FINAL RESULTS                ║")
	fmt.Println("╠════════════════════════════════════════╣")
	
	if stats["has_winner"].(bool) {
		fmt.Printf("║ 🏆 WINNER: User %-22d ║\n", stats["winner_user_id"])
		fmt.Printf("║    Answer: %-27s ║\n", stats["winner_answer"])
		timeToWin := stats["time_to_win"].(float64)
		fmt.Printf("║    Time to win: %.3f seconds         ║\n", timeToWin)
	} else {
		fmt.Println("║ ❌ No winner found (no correct answers) ║")
	}
	
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Println("║            STATISTICS                 ║")
	fmt.Println("╠════════════════════════════════════════╣")
	
	total := stats["total_responses"].(int64)
	correct := stats["correct_responses"].(int64)
	
	fmt.Printf("║ Total Responses: %-21d ║\n", total)
	fmt.Printf("║ Correct Responses: %-19d ║\n", correct)
	
	if total > 0 {
		percentage := stats["correct_percentage"].(float64)
		fmt.Printf("║ Success Rate: %.2f%%                   ║\n", percentage)
	}
	
	fmt.Printf("║ Total Time: %.3f seconds              ║\n", time.Since(startTime).Seconds())
	fmt.Println("╚════════════════════════════════════════╝")
}

func handleShutdown(engine *game_engine.GameEngine) {
	fmt.Println("\n🛑 Shutting down server...")
	
	if winner := engine.GetWinner(); winner != nil {
		fmt.Printf("Final Winner: User %d with answer '%s'\n", winner.UserID, winner.Answer)
	}
	
	stats := engine.GetStats()
	fmt.Printf("Total responses processed: %d\n", stats["total_responses"])
	
	engine.Shutdown()
	os.Exit(0)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func printBanner(title string) {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Printf("║         %-32s ║\n", title)
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
}