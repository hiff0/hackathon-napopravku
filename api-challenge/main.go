package main

import (
	"api-challenge/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

var log = setupPrettySlog()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ConnectionMessage struct {
	UsersCount int `json:"usersCount"`
}

type server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan int
	mu        sync.Mutex
}

func newServer() *server {
	return &server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan int),
	}
}

func (s *server) handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Upgrade to WebSocket error:", err.Error())
		return
	}
	defer ws.Close()

	s.mu.Lock()
	s.clients[ws] = true
	totalConnections := len(s.clients)
	s.mu.Unlock()

	s.broadcast <- totalConnections

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			s.mu.Lock()
			delete(s.clients, ws)
			totalConnections = len(s.clients)
			s.mu.Unlock()
			s.broadcast <- totalConnections
			break
		}
	}
}

func (s *server) handleBroadcast() {
	for {
		totalConnections := <-s.broadcast
		msg := ConnectionMessage{
			UsersCount: totalConnections,
		}

		s.mu.Lock()
		for client := range s.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Error("WriteJSON error:", err.Error())
				client.Close()
				delete(s.clients, client)
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	wsServer := newServer()

	http.HandleFunc("/ws", wsServer.handleConnections)

	go wsServer.handleBroadcast()

	log.Info("Server started on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Error("Server start error:", err.Error())
	}
}
