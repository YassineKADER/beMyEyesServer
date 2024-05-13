package utils

import (
	"os"
	"testing"
)

func TestSts(t *testing.T) {
	data, err := os.ReadFile("../../testdata/testAudio.wav")
	if err != nil {
		t.Fatal(err)
	}
	sentence := SpeakToText(data)
	if sentence == "" {
		t.Errorf("sentence is empty")
	}
}
