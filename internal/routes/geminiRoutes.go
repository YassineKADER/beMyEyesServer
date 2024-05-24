package routes

import (
	"github.com/YassineKADER/beMyEyesServer/internal/handlers"
	imagenet "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/YassineKADER/beMyEyesServer/internal/utils"
	"github.com/go-chi/chi/v5"
)

func geminiRoute(model *imagenet.Model, ocrInstance *ocr.OCR, gemini *utils.GeminiModel, geminiVision *utils.GeminiModel) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", handlers.GeminiHandler(model, ocrInstance, gemini))
	router.Post("/vision", handlers.GeminiVisionHandler(geminiVision))
	router.Post("/audio", handlers.StsHandler(gemini))
	return router
}
