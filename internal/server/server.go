package server

import (
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/config"
	"github.com/YassineKADER/beMyEyesServer/internal/routes"
)

func NewServer() error {
	var conf config.Config
	conf.Load()
	router := routes.NewRouter()
	return http.ListenAndServe("localhost:"+conf.Port, router)
}
