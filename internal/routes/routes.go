package routes

import (
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/config"
	imagenet "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/go-chi/chi/v5"
)

func Routes(imagenetModel *imagenet.Model, ocrModel *ocr.OCR) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/imagenet", imagenetRoute(imagenetModel))
		r.Mount("/api/ocr", ocrRoute(ocrModel))
	})
	return router
}

func CreateRouter(modeldir string) *chi.Mux {
	model := imagenet.Model{}
	ocr := ocr.OCR{}
	if modeldir == "" {
		modeldir = "./modeldir"
	}
	model.Load(modeldir)
	ocr.Load()
	var conf config.Config
	conf.Load()
	router := Routes(&model, &ocr)
	// defer ocr.Close()
	// defer model.Close()
	return router
}
