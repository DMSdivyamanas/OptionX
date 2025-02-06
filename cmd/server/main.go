// cmd/server/main.go
package main

import (
	"OptionX_Assignment/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// the entry point of the application.
func main() {
	// Create a new server instance.
	srv := server.NewServer()

	// Start the server in a separate goroutine.
	go srv.Run()

	// Set up the WebSocket endpoint.
	http.HandleFunc("/ws", srv.HandleConnections)

	// Handle OS signals for graceful shutdown.
	go func() {
		// Create a channel to listen for OS signals.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		// Block until a signal is received.
		<-c
		log.Println("Shutting down server...")
		// Shutdown the server gracefully and exit
		srv.Shutdown()
		os.Exit(0)
	}()

	// Logging that that the server has started.
	log.Println("Server started on :8081")
	// Start the HTTP server on port 8081.
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}