package routes

import (
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/config"
	imagenet "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/YassineKADER/beMyEyesServer/internal/middleware"
	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/YassineKADER/beMyEyesServer/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(imagenetModel *imagenet.Model, ocrModel *ocr.OCR, gemini *utils.GeminiModel) *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "User-Agent"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(middleware.LoggingMiddleware)
	router.Use(cors.Handler)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/imagenet", imagenetRoute(imagenetModel))
		r.Mount("/api/ocr", ocrRoute(ocrModel))
		r.Mount("/api/gemini", geminiRoute(imagenetModel, ocrModel, gemini))
	})
	return router
}

func CreateRouter(modeldir string) *chi.Mux {
	model := imagenet.Model{}
	ocr := ocr.OCR{}
	gemini := utils.GeminiModel{}
	if modeldir == "" {
		modeldir = "./modeldir"
	}
	model.Load(modeldir)
	ocr.Load()
	gemini.LoadModel("models/gemini-pro")
	var conf config.Config
	conf.Load()
	router := Routes(&model, &ocr, &gemini)
	// defer ocr.Close()
	// defer model.Close()
	return router
}
