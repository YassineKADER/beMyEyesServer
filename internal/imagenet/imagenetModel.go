package imagenetModel

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
)

type Model struct {
	session                *tf.Session
	graph                  *tf.Graph
	imageNormalizerGraph   *tf.Graph
	imageNormalizerSession *tf.Session
	imageNormalizerInput   tf.Input
	imageNormalizerOutput  tf.Output
	lablelStrings          []string
}

func (m *Model) Load(modeldir string) error {
	modelfile, labelsfile, err := modelFiles(modeldir)
	if err != nil {
		log.Fatal(err)
	}
	model, err := os.ReadFile(modelfile)
	if err != nil {
		log.Fatal(err)
	}
	m.graph = tf.NewGraph()
	if err := m.graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}
	m.session, err = tf.NewSession(m.graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(labelsfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m.lablelStrings = append(m.lablelStrings, scanner.Text())
	}
	m.imageNormalizerGraph, m.imageNormalizerInput, m.imageNormalizerOutput, err = constructGraphToNormalizeImage()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (m *Model) Close() {
	m.session.Close()
	m.imageNormalizerSession.Close()
}

func modelFiles(dir string) (modelfile, labelsfile string, err error) {
	const URL = "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip"
	var (
		model   = filepath.Join(dir, "tensorflow_inception_graph.pb")
		labels  = filepath.Join(dir, "imagenet_comp_graph_label_strings.txt")
		zipfile = filepath.Join(dir, "inception5h.zip")
	)
	if filesExist(model, labels) == nil {
		return model, labels, nil
	}
	log.Println("Did not find model in", dir, "downloading from", URL)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", err
	}
	if err := download(URL, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to download %v - %v", URL, err)
	}
	if err := unzip(dir, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to extract contents from model archive: %v", err)
	}
	os.Remove(zipfile)
	return model, labels, filesExist(model, labels)
}

func filesExist(files ...string) error {
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("unable to stat %s: %v", f, err)
		}
	}
	return nil
}

func download(URL, filename string) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func unzip(dir, zipfile string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		src, err := f.Open()
		if err != nil {
			return err
		}
		log.Println("Extracting", f.Name)
		dst, err := os.OpenFile(filepath.Join(dir, f.Name), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		dst.Close()
	}
	return nil
}

func constructGraphToNormalizeImage() (graph *tf.Graph, input, output tf.Output, err error) {
	const (
		H, W  = 224, 224
		Mean  = float32(117)
		Scale = float32(1)
	)
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s,
						op.DecodeJpeg(s, input, op.DecodeJpegChannels(3)), tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}
