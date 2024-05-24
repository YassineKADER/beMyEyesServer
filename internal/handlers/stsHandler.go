package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/YassineKADER/beMyEyesServer/internal/utils"
	"github.com/go-chi/render"
)

func StsHandler(gemini *utils.GeminiModel) http.HandlerFunc {
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
			file, _, err := r.FormFile("audio")
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

		result := utils.SpeakToText(body)
		fmt.Println(result)
		response := gemini.GenerateResponseForQuestion(result)
		fmt.Println(response)
		responseMap := map[string]string{"response": response}
		jsonData, err := json.Marshal(responseMap)
		if err != nil {
			render.Status(r, 500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
