package main

import (
	"log"

	"github.com/YassineKADER/beMyEyesServer/internal/server"
)

func main() {
	if err := server.NewServer(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
