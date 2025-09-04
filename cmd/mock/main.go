package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/glitchdawg/game-engine-with-user/mock_engine"
)

func main() {
	var numUsers int
	var apiURL string

	flag.IntVar(&numUsers, "users", 100, "Number of users to simulate")
	flag.StringVar(&apiURL, "api", "http://localhost:8080/submit", "API server URL")
	flag.Parse()

	rand.NewSource(45)//RANDOM SEED GENERATOR

	fmt.Printf("Mock User Engine Starting\n")
	fmt.Printf("Number of users: %d\n", numUsers)
	fmt.Printf("API URL: %s\n\n", apiURL)

	engine := mock_engine.NewMockEngine(apiURL)
	
	start := time.Now()
	engine.SimulateUsers(numUsers)
	
	fmt.Printf("\nSimulation completed in %v\n", time.Since(start))
	log.Println("Mock engine finished")
}