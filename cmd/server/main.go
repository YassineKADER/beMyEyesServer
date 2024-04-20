package main

import (
	"fmt"
	"time"

	imagenetModel "github.com/YassineKADER/beMyEyesServer/internal/imagenet"
)

func main() {
	model := imagenetModel.Model{}
	startLoad := time.Now()
	model.Load("modeldir")
	loadTime := time.Since(startLoad)
	fmt.Printf("Time to load model: %s\n", loadTime)
	startMatch := time.Now()
	label := model.Match("test.jpg", false)
	matchTime := time.Since(startMatch)
	fmt.Printf("Time to match: %s\n", matchTime)
	fmt.Println(label)
}
