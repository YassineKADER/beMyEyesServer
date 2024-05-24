package routes

import (
	"github.com/YassineKADER/beMyEyesServer/internal/handlers"
	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/go-chi/chi/v5"
)

func ocrRoute(ocrInstance *ocr.OCR) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", handlers.OcrHandler(ocrInstance))
	return router
}
