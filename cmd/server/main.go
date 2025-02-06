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

func main() {
	srv := server.NewServer()
	go srv.Run()

	http.HandleFunc("/ws", srv.HandleConnections)

	// Handle OS signals for graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")
		srv.Shutdown()
		os.Exit(0)
	}()

	log.Println("Server started on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}