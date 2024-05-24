package imagenetModel

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	tf "github.com/wamuir/graft/tensorflow" // graft
	"github.com/wamuir/graft/tensorflow/op"
)

type Model struct {
	session                *tf.Session
	graph                  *tf.Graph
	imageNormalizerGraph   *tf.Graph
	imageNormalizerSession *tf.Session
	imageNormalizerInput   tf.Output
	imageNormalizerOutput  tf.Output
	lablelStrings          []string
}

type Match struct {
	Label       string  `json:"label"`
	Number      int     `json:"number"`
	Probability float32 `json:"probability"`
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
	m.imageNormalizerSession, err = tf.NewSession(m.imageNormalizerGraph, nil)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (m *Model) Match(imagefile string, url bool, data *[]byte) []byte {
	tensor, err := m.makeTensorFromImage(imagefile, url, data)
	if err != nil {
		log.Fatal(err)
	}
	// Construct an in-memory graph from the serialized form.
	probabilities, err := m.session.Run(
		map[tf.Output]*tf.Tensor{
			m.graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			m.graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	return printBestLabel(probabilities[0].Value().([][]float32)[0], *m)
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

// func printBestLabel(probabilities []float32, m Model) string {
// 	bestIdxs := make([]int, 3)
// 	copy(bestIdxs, []int{0, 0, 0})
//
// 	for i, p := range probabilities {
// 		for j := 0; j < 3; j++ {
// 			if p > probabilities[bestIdxs[j]] {
// 				copy(bestIdxs[j+1:], bestIdxs[j:])
// 				bestIdxs[j] = i
// 				break
// 			}
// 		}
// 	}
//
// 	result := ""
// 	for i, idx := range bestIdxs {
// 		result += fmt.Sprintf("MATCH %d: (%2.0f%% likely) %s\n", i+1, probabilities[idx]*100.0, m.lablelStrings[idx])
// 	}
//
// 	return result
// }

func printBestLabel(probabilities []float32, m Model) []byte {
	bestIdxs := make([]int, 3)
	copy(bestIdxs, []int{0, 0, 0})

	for i, p := range probabilities {
		for j := 0; j < 3; j++ {
			if p > probabilities[bestIdxs[j]] {
				copy(bestIdxs[j+1:], bestIdxs[j:])
				bestIdxs[j] = i
				break
			}
		}
	}

	matches := make([]Match, 3)
	for i, idx := range bestIdxs {
		matches[i] = Match{
			Number:      i + 1,
			Probability: probabilities[idx] * 100.0,
			Label:       m.lablelStrings[idx],
		}
	}

	jsonData, err := json.Marshal(matches)
	if err != nil {
		return []byte("error")
	}

	return jsonData
}

// Convert the image in filename to a Tensor suitable as input to the Inception model.
func (m *Model) makeTensorFromImage(filename string, url bool, data *[]byte) (*tf.Tensor, error) {
	var bytes []byte
	if data != nil && len(*data) > 0 {
		bytes = *data
	} else {
		if url {
			res, err := http.Get(filename)
			if err != nil {
				return nil, errors.New("invalid image url")
			}
			defer res.Body.Close()
			bytes, err = io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.New("invalid image url")
			}
		} else {
			var err error
			bytes, err = os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
		}
	}
	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	// Execute that graph to normalize this one image
	normalized, err := m.imageNormalizerSession.Run(
		map[tf.Output]*tf.Tensor{m.imageNormalizerInput: tensor},
		[]tf.Output{m.imageNormalizerOutput},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}
