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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn)
	s.register <- client

	go s.handlePingPong(client)

	defer func() {
		s.unregister <- client
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("Received from client %s: %s", client.ID, message)

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		if targetID, ok := msg["id"].(string); ok {
			s.mu.Lock()
			if targetID == "broadcast" {
				broadcastMessage := fmt.Sprintf("%s: %s", client.Username, msg["message"])
				for _, client := range s.clients {
					err := client.Conn.WriteMessage(websocket.TextMessage, []byte(broadcastMessage))
					if err != nil {
						log.Printf("Error broadcasting to client %s: %v", client.ID, err)
					}
				}
			} else if targetClient, exists := s.clients[targetID]; exists {
				directMessage := fmt.Sprintf("Direct message from %s: %s", client.Username, msg["message"])
				err := targetClient.Conn.WriteMessage(websocket.TextMessage, []byte(directMessage))
				if err != nil {
					log.Printf("Error sending to client %s: %v", targetClient.ID, err)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) handlePingPong(client *Client) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	client.Conn.SetPongHandler(func(appData string) error {
		// log.Printf("Received pong from client %s", client.ID)
		return nil
	})

	for range ticker.C {
		if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Ping error: %v", err)
			s.unregister <- client
			return
		}
	}
}