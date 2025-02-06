// internal/server/server.go
package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()
			welcomeMessage := "Welcome! Current clients: " + s.getClientIDs()
			client.Conn.WriteMessage(websocket.TextMessage, []byte(welcomeMessage))

			newClientMessage := fmt.Sprintf("New client connected: %s", client.ID)
			s.broadcastMessageToAll([]byte(newClientMessage), client.ID)
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				delete(s.clients, client.ID)
				closeConnection(client.Conn)
			}
			s.mu.Unlock()
		case message := <-s.broadcast:
			s.mu.Lock()
			for _, client := range s.clients {
				client.Conn.WriteMessage(websocket.TextMessage, message)
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) broadcastMessageToAll(message []byte, excludeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, client := range s.clients {
		if id != excludeID {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Printf("Error broadcasting to client %s: %v\n", client.ID, err)
			}
		}
	}
}

func (s *Server) getClientIDs() string {
	ids := ""
	for id := range s.clients {
		ids += id + " "
	}
	return ids
}

func (s *Server) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		closeConnection(client.Conn)
	}
}

func closeConnection(conn *websocket.Conn) {
	conn.Close()
}