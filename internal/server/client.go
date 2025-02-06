// internal/server/client.go
package server

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:       uuid.New().String(),
		Username: gofakeit.Username(),
		Conn:     conn,
	}
}