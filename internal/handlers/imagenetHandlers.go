package handlers

import (
	"io"
	"net/http"
	"strings"

	imagenet "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
	"github.com/go-chi/render"
)

// func ImagenetHandler(model *imagenet.Model) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			render.Status(r, 500)
// 			return
// 		}
// 		result := model.Match("", false, &body)
// 		render.JSON(w, r, result)
// 	}
// }

func ImagenetHandler(model *imagenet.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		var err error

		// Check the Content-Type of the request
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			// Parse the multipart form in the request
			err = r.ParseMultipartForm(10 << 20) // Maximum memory size is 10MB
			if err != nil {
				render.Status(r, 500)
				return
			}

			// Get a reference to the file field in the form
			file, _, err := r.FormFile("image")
			if err != nil {
				render.Status(r, 500)
				return
			}
			defer file.Close()

			// Read the file data
			body, err = io.ReadAll(file)
			if err != nil {
				render.Status(r, 500)
				return
			}
		} else {
			// Read the body data
			body, err = io.ReadAll(r.Body)
			if err != nil {
				render.Status(r, 500)
				return
			}
		}

		result := model.Match("", false, &body)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		// render.JSON(w, r, result)
	}
}
