package handlers

import (
	"io"
	"net/http"

	imagenet "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/go-chi/render"
)

func ImagenetHandler(model *imagenet.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.Status(r, 500)
			return
		}
		result := model.Match("", false, &body)
		render.JSON(w, r, result)
	}
}
