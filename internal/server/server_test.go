package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/YassineKADER/beMyEyesServer/internal/routes"
)

func TestServerImagenetApi(t *testing.T) {
	router := routes.CreateRouter("./../../modeldir")
	ts := httptest.NewServer(router)
	defer ts.Close()
	imgData, err := os.ReadFile("./../../testdata/test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/v1/api/imagenet", bytes.NewBuffer(imgData))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestServerOcrapi(t *testing.T) {
	router := routes.CreateRouter("./../../modeldir")
	ts := httptest.NewServer(router)
	defer ts.Close()
	imgData, err := os.ReadFile("./../../testdata/test_ocr.jpg")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", ts.URL+"/v1/api/ocr", bytes.NewBuffer(imgData))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
