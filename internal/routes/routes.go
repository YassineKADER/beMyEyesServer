package routes

import (
	"net/http"

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
