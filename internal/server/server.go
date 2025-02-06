// internal/server/server.go
package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Server struct holds information about connected clients and channels for communication.
type Server struct {
	clients    map[string]*Client // Map of client IDs to Client objects
	register   chan *Client       // Channel for registering new clients
	unregister chan *Client       // Channel for unregistering clients
	broadcast  chan []byte        // Channel for broadcasting messages to all clients
	mu         sync.Mutex         // Mutex to protect access to the clients map
}

// NewServer initializes a new Server instance.
func NewServer() *Server {
	return &Server{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run starts the server's main loop to handle client registration, unregistration, and broadcasting.
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			// Register a new client
			s.mu.Lock()
			s.clients[client.ID] = client
			s.mu.Unlock()

			// Send a welcome message to the new client
			welcomeMessage := "Welcome! Current clients: " + s.getClientUsernames()
			client.Conn.WriteMessage(websocket.TextMessage, []byte(welcomeMessage))

			// Broadcast the new client's username to all other clients
			newClientMessage := fmt.Sprintf("New client connected: %s", client.Username)
			s.broadcastMessageToAll([]byte(newClientMessage), client.ID)

		case client := <-s.unregister:
			// Unregister a client
			s.mu.Lock()
			if _, ok := s.clients[client.ID]; ok {
				delete(s.clients, client.ID)
				closeConnection(client.Conn)
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			// Broadcast a message to all clients
			s.mu.Lock()
			for _, client := range s.clients {
				client.Conn.WriteMessage(websocket.TextMessage, message)
			}
			s.mu.Unlock()
		}
	}
}

// broadcastMessageToAll sends a message to all clients except the one specified by excludeID.
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

// getClientUsernames returns a string of all connected client usernames.
func (s *Server) getClientUsernames() string {
	usernames := ""
	for _, client := range s.clients {
		usernames += client.Username + " "
	}
	return usernames
}

// Shutdown gracefully shuts down the server by closing all client connections.
func (s *Server) Shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, client := range s.clients {
		closeConnection(client.Conn)
	}
}

// closeConnection closes a WebSocket connection.
func closeConnection(conn *websocket.Conn) {
	conn.Close()
}