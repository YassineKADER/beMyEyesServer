package ocr

import (
	"os"
	"testing"
)

func TestOcrLoad(t *testing.T) {
	o := OCR{}
	o.Load()
	defer o.Close()
	if o.client == nil {
		t.Errorf("o.client is nil")
	}
}

func TestOcrRecognize(t *testing.T) {
	o := OCR{}
	o.Load()
	defer o.Close()
	text, err := o.Recognize("./../../testdata/test_ocr.jpg")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if text == "" {
		t.Errorf("text is empty")
	}
}

func TestOcrRecognizeFromBytes(t *testing.T) {
	o := OCR{}
	o.Load()
	defer o.Close()
	bytes, _ := os.ReadFile("./../../testdata/test_ocr.jpg")
	text, err := o.RecognizeFromBytes(bytes)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if text == "" {
		t.Errorf("text is empty")
	}
}
