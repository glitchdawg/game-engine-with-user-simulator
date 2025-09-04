package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/glitchdawg/game-engine-with-user/api_server"
	"github.com/glitchdawg/game-engine-with-user/game_engine"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Server port")
	flag.Parse()

	fmt.Println("===========================================")
	fmt.Println("       Game API Server Starting")
	fmt.Println("===========================================")
	fmt.Printf("Port: %s\n", port)
	fmt.Println()

	engine := game_engine.NewGameEngine()
	
	server := api_server.NewAPIServer(port, engine)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		
		if winner := engine.GetWinner(); winner != nil {
			fmt.Printf("\nFinal Winner: User %d with answer '%s'\n", winner.UserID, winner.Answer)
		}
		os.Exit(0)
	}()

	fmt.Println("Server is ready to receive requests")
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  POST /submit - Submit user responses\n")
	fmt.Printf("  GET  /stats  - View current statistics\n")
	fmt.Printf("  POST /reset  - Reset the game\n")
	fmt.Println("\nPress Ctrl+C to stop the server")
	fmt.Println("-------------------------------------------")

	if err := server.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}