package ocr

import (
	"io"
	"log"
	"net/http"

	"github.com/otiai10/gosseract"
)

type ocr struct {
	client *gosseract.Client
}

func (o *ocr) Load() {
	o.client = gosseract.NewClient()
}

func (o *ocr) Close() {
	o.client.Close()
}

func (o *ocr) Recognize(imagePath string) (string, error) {
	o.client.SetImage(imagePath)
	return o.client.Text()
}

func (o *ocr) RecognizeFromBytes(imageBytes []byte) (string, error) {
	o.client.SetImageFromBytes(imageBytes)
	return o.client.Text()
}

func (o *ocr) RecognizeFromURL(imageURL string) (string, error) {
	if imageURL != "" {
		res, err := http.Get(imageURL)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		o.client.SetImageFromBytes(bytes)
		return o.client.Text()
	}
	return "", nil
}
