package utils

import (
	"context"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

func SpeakToText(data []byte) string {
	context := context.Background()
	client, err := speech.NewClient(context)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return ""
	}
	resp, err := client.Recognize(context, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                            speechpb.RecognitionConfig_MP3,
			LanguageCode:                        "en-US",
			AudioChannelCount:                   2,
			EnableSeparateRecognitionPerChannel: false,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
		return ""
	}
	results := ""
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			results += alt.Transcript
		}
	}
	return results
}
