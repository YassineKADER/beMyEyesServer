package imagenetModel

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tf "github.com/wamuir/graft/tensorflow" // graft
	// graft
)

// Match, modified from:
// https://pkg.go.dev/github.com/tensorflow/tensorflow/tensorflow/go#example-package
// An example for using the TensorFlow Go API for image recognition
// using a pre-trained inception model (http://arxiv.org/abs/1512.00567).
//
// tesmgnetModel.got// Sample usage: <program> -dir=/tmp/modeldir -image=/path/to/some/jpeg
//
// The pre-trained model takes input in the form of a 4-dimensional
// tensor with shape [ BATCH_SIZE, IMAGE_HEIGHT, IMAGE_WIDTH, 3 ],
// where:
// - BATCH_SIZE allows for inference of multiple images in one pass through the graph
// - IMAGE_HEIGHT is the height of the images on which the model was trained
// - IMAGE_WIDTH is the width of the images on which the model was trained
// - 3 is the (R, G, B) values of the pixel colors represented as a float.
//
// And produces as output a vector with shape [ NUM_LABELS ].
// output[i] is the probability that the input image was recognized as
// having the i-th label.
//
// A separate file contains a list of string labels corresponding to the
// integer indices of the output.
//
// This example:
// - Loads the serialized representation of the pre-trained model into a Graph
// - Creates a Session to execute operations on the Graph
// - Converts an image file to a Tensor to provide as input to a Session run
// - Executes the Session and prints out the label with the highest probability
//
// To convert an image file to a Tensor suitable for input to the Inception model,
// this example:
//   - Constructs another TensorFlow graph to normalize the image into a
//     form suitable for the model (for example, resizing the image)
//   - Creates and executes a Session to obtain a Tensor in this normalized form.
func Match(modeldir, imagefile *string, url bool) string {
	modelfile, labelsfile, err := modelFiles(*modeldir)
	if err != nil {
		log.Fatal(err)
	}
	model, err := ioutil.ReadFile(modelfile)
	if err != nil {
		log.Fatal(err)
	}

	// Construct an in-memory graph from the serialized form.
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	// Create a session for inference over graph.
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Run inference on *imageFile.
	// For multiple images, session.Run() can be called in a loop (and
	// concurrently). Alternatively, images can be batched since the model
	// accepts batches of image data as input.
	tensor, err := makeTensorFromImage(*imagefile, url)
	if err != nil {
		log.Fatal(err)
	}
	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	// output[0].Value() is a vector containing probabilities of
	// labels for each image in the "batch". The batch size was 1.
	// Find the most probably label index.
	probabilities := output[0].Value().([][]float32)[0]
	return printBestLabel(probabilities, labelsfile)
}

func printBestLabel(probabilities []float32, labelsFile string) string {
	bestIdx := 0
	for i, p := range probabilities {
		if p > probabilities[bestIdx] {
			bestIdx = i
		}
	}
	// Found the best match. Read the string from labelsFile, which
	// contains one line per label.
	file, err := os.Open(labelsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ERROR: failed to read %s: %v", labelsFile, err)
	}
	return fmt.Sprintf("BEST MATCH: (%2.0f%% likely) %s\n", probabilities[bestIdx]*100.0, labels[bestIdx])
}

// Convert the image in filename to a Tensor suitable as input to the Inception model.
func makeTensorFromImage(filename string, url bool) (*tf.Tensor, error) {
	var bytes []byte
	if url {
		res, err := http.Get(filename)
		if err != nil {
			return nil, errors.New("invalid image url")
		}
		defer res.Body.Close()
		bytes, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.New("invalid image url")
		}
	} else {
		var err error
		bytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	// DecodeJpeg uses a scalar String-valued tensor as input.
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}
	// Construct a graph to normalize the image
	graph, input, output, err := constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}
	// Execute that graph to normalize this one image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as
// input a JPEG-encoded string and returns a tensor suitable as input to the
// inception model.
