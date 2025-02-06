// internal/server/handlers.go
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleConnections manages incoming WebSocket connections.
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client with the WebSocket connection.
	client := NewClient(conn)
	// Register the new client with the server.
	s.register <- client

	// Start handling ping/pong messages for the client.
	go s.handlePingPong(client)

	// Ensure the client is unregistered when the function exits.
	defer func() {
		s.unregister <- client
	}()

	// Continuously read messages from the client.
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Received from client %s: %s", client.ID, message)

		// Parse the incoming message as JSON.
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		// Determine the target client for the message.
		if targetID, ok := msg["id"].(string); ok {
			s.mu.Lock()
			if targetID == "broadcast" {
				// Broadcast the message to all clients.
				broadcastMessage := fmt.Sprintf("%s: %s", client.Username, msg["message"])
				for _, client := range s.clients {
					client.Conn.WriteMessage(websocket.TextMessage, []byte(broadcastMessage))
				}
			} else if targetClient, exists := s.clients[targetID]; exists {
				// Send a direct message to the specified client.
				directMessage := fmt.Sprintf("Direct message from %s: %s", client.Username, msg["message"])
				targetClient.Conn.WriteMessage(websocket.TextMessage, []byte(directMessage))
			}
			s.mu.Unlock()
		}
	}
}

// handlePingPong manages the ping/pong mechanism for a client.
func (s *Server) handlePingPong(client *Client) {
	// Create a ticker to send ping messages at regular intervals.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Set a handler for pong messages from the client.
	client.Conn.SetPongHandler(func(appData string) error {
		log.Printf("Received pong from client %s", client.ID)
		return nil
	})

	// Continuously send ping messages.
	for range ticker.C {
		log.Printf("Sending ping to client %s", client.ID)
		if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Ping error: %v", err)
			s.mu.Lock()
			// Broadcast a disconnection message to all other clients.
			disconnectionMessage := fmt.Sprintf("%s is disconnected", client.Username)
			for _, c := range s.clients {
				if c.ID != client.ID {
					c.Conn.WriteMessage(websocket.TextMessage, []byte(disconnectionMessage))
				}
			}
			// Unregister and delete the client.
			delete(s.clients, client.ID)
			s.mu.Unlock()
			return
		}
	}
}