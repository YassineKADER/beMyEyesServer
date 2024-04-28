package handlers

import (
	"io"
	"net/http"

	"github.com/YassineKADER/beMyEyesServer/internal/ocr"
	"github.com/go-chi/render"
)

func OcrHandler(ocrInstance *ocr.OCR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.Status(r, 500)
			return
		}
		result, err := ocrInstance.RecognizeFromBytes(body)
		if err != nil {
			render.Status(r, 500)
			return
		}
		render.JSON(w, r, result)
	}
}
