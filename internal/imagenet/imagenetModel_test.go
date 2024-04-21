package imagenetModel

import (
	"fmt"
	"testing"
)

func TestImagenetModelLoad(t *testing.T) {
	model := Model{}
	model.Load("modeldir")
	if model.session == nil {
		t.Errorf("model.session is nil")
	}
	if model.graph == nil {
		t.Errorf("model.graph is nil")
	}
	if model.imageNormalizerGraph == nil {
		t.Errorf("model.imageNormalizerGraph is nil")
	}
	if model.imageNormalizerSession == nil {
		t.Errorf("model.imageNormalizerSession is nil")
	}
}

func TestImagenetModelPredict(t *testing.T) {
	model := Model{}
	model.Load("modeldir")
	prediction := model.Match("https://i.ibb.co/Y2s0WH6/test-dog.jpg", true)
	if prediction == "" {
		t.Errorf("prediction is nil")
	}
	fmt.Println(prediction)
}

func TestImagenetModelClose(t *testing.T) {
	model := Model{}
	model.Load("modeldir")
	model.Close()
}
