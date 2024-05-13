package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiModel struct {
	context context.Context
	client  *genai.Client
	model   *genai.GenerativeModel
}

func (gm *GeminiModel) LoadModel(name string) {
	gm.context = context.Background()
	var err error
	gm.client, err = genai.NewClient(gm.context, option.WithAPIKey(os.Getenv("GEMINI_KEY")))
	if err != nil {
		return
	}
	gm.model = gm.client.GenerativeModel(name)
}

func (gm *GeminiModel) GenerateResponse(prompt string) string {
	prompt = "Your task is to identify the most likely object and describe it in simple terms please combain the 3 objects, and the data from ocr text and find commun thing and tell me what this object can be, if the text was too long just talk about this text and what it's content. For instance, Keep the description short and clear, suitable for helping blind individuals understand what the object might be" + prompt + "\n"
	output, err := gm.model.GenerateContent(gm.context, genai.Text(prompt))
	if err != nil {
		return ""
	}
	return formatResponse(output)
}

func (gm *GeminiModel) GenerateResponseFromPicture(img []byte) string {
	prompt := []genai.Part{
		genai.ImageData("jpg", img),
		genai.Text("describe what in the picture for guid a blind user of our app"),
	}
	output, err := gm.model.GenerateContent(gm.context, prompt...)
	if err != nil {
		return ""
	}
	return formatResponse(output)
}

func formatResponse(resp *genai.GenerateContentResponse) string {
	var formattedContent strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, cand := range resp.Candidates {
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					formattedContent.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	return formattedContent.String()
}
