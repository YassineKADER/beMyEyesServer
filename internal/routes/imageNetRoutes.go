package routes

import (
	"github.com/YassineKADER/beMyEyesServer/internal/handlers"
	imagenetModel "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/go-chi/chi/v5"
)

func imagenetRoute(Model *imagenetModel.Model) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", handlers.ImagenetHandler(Model))
	return router
}
