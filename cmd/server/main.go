package main

import "github.com/YassineKADER/beMyEyesServer/internal/server"

func main() {
	err := server.NewServer()
	if err != nil {
		panic(err)
	}
}
