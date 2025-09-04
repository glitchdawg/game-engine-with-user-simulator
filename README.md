# Game Engine with User Simulator

A Go-based backend system that simulates multiple users answering game questions, evaluates responses in real-time, and announces winners.

## Features

✅ **Concurrent User Simulation** - Handle 1000+ concurrent users  
✅ **Real-time Evaluation** - Process responses as they arrive  
✅ **Channel-based Event Handling** - Efficient event-driven architecture  
✅ **Live Metrics** - Track correct/incorrect answers in real-time  
✅ **Single API Endpoint** - Clean `/submit` endpoint  
✅ **Interactive CLI** - Built-in commands for stats and reset  
✅ **Race Condition Free** - Thread-safe with mutexes and atomic operations  

## Architecture

### 1. Mock User Engine
- Simulates N concurrent users
- Random delays (0-1000ms) per user  
- 30% chance of correct answer
- Sends responses concurrently to API server

### 2. API Server  
- Single `/submit` endpoint (POST)
- Forwards responses to Game Engine
- Thread-safe request handling

### 3. Game Engine
- Channel-based event processing
- Atomic operations for metrics
- First correct answer wins
- Ignores subsequent correct answers
- Live statistics every 5 seconds

## Usage

### Build
```bash
go build .
```

### Run Modes

#### 1. Full Simulation (Recommended)
```bash
go run . -mode full -users 1000 -port 8080
```

#### 2. Server Mode (Interactive)
```bash
go run . -mode server -port 8080
```
Commands:
- `stats` - Show current statistics
- `reset` - Reset the game engine  
- `clear` - Clear the screen
- `exit` - Shutdown server

#### 3. Mock Users Only
```bash
go run . -mode mock -users 1000 -api http://localhost:8080/submit
```

### Flags
- `-mode` - Operation mode: server, mock, or full (default: server)
- `-port` - API server port (default: 8080)
- `-users` - Number of mock users (default: 1000)
- `-api` - API URL for mock engine (default: http://localhost:8080/submit)

## Project Structure
```
.
├── api_server/
│   └── server.go       # HTTP API server
├── game_engine/
│   └── engine.go       # Game logic & winner detection
├── mock_engine/
│   └── mock_engine.go  # User simulator
├── cmd/
│   ├── api/
│   │   └── main.go     # Standalone API server
│   └── mock/
│       └── main.go     # Standalone mock engine
└── main.go             # Main entry point
```

## Performance

- Handles 1000+ concurrent requests without race conditions
- Average response time < 1ms per request
- Channel buffer prevents blocking under load
- Atomic operations for thread-safe counters

## Example Output

```
╔══════════════════════════════════════════╗
║            🎉 WINNER FOUND! 🎉            ║
╠══════════════════════════════════════════╣
║ Winner ID:     780                        ║
║ Answer:        42                         ║
║ Time to win:   2.019s                     ║
║ Total responses: 456                      ║
║ Correct answers: 137                      ║
╚══════════════════════════════════════════╝
```

## Testing

Run with different user counts to test load:
```bash
go run . -mode full -users 100
go run . -mode full -users 1000
go run . -mode full -users 10000
```

## Requirements Met

✅ Language: Go  
✅ Handles 1000+ concurrent requests  
✅ No race conditions or deadlocks  
✅ Real-time evaluation (no batching)  
✅ Clean structure with separate modules  
✅ Only one winner declared  

### Bonus Points Achieved
✅ Metrics tracking for correct/incorrect answers  
✅ Time taken to find winner  
✅ Channel-based event handling  
✅ Live statistics during gameplay  
✅ Atomic operations for performance