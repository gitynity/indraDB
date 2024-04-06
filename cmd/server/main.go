package main

import (
	"log"

	"github.com/gitynity/indraDB/pkg/server"
)

func main() {

	// Create a new instance of the server
	s := server.NewServer(":8080", "/Users/nns/indraDB")

	// Start the server
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
