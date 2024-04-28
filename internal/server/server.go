package server

import (
	"fmt"
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/config"
	imagenetModel "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/YassineKADER/beMyEyesServer/internal/routes"
)

func NewServer() error {
	model := imagenetModel.Model{}
	ocr := ocr.OCR{}
	model.Load("modeldir")
	ocr.Load()
	var conf config.Config
	conf.Load()
	router := routes.Routes(&model, &ocr)
	fmt.Println("Listening on localhost:" + conf.Port)
	return http.ListenAndServe("localhost:"+conf.Port, router)
}
