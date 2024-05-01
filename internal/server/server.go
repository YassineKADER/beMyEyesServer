package server

import (
	"fmt"
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/config"
	"github.com/YassineKADER/beMyEyesServer/internal/routes"
)

func NewServer() error {
	var conf config.Config
	conf.Load()
	router := routes.CreateRouter("")
	fmt.Println("Listening on localhost:" + conf.Port)
	return http.ListenAndServe("0.0.0.0:"+conf.Port, router)
}
