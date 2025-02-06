
package server

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client struct represents a connected client with a unique ID and username.
type Client struct {
	ID       string          // Unique identifier for the client
	Username string          // Randomly generated username for the client
	Conn     *websocket.Conn // WebSocket connection for the client
}

// NewClient creates a new Client instance with a WebSocket connection.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:       uuid.New().String(), // Generate a unique UUID for the client
		Username: gofakeit.Username(), // Generate a random username for the client
		Conn:     conn,                // Assign the WebSocket connection
	}
}