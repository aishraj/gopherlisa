package main

import (
	"flag"
	"github.com/aishraj/gopherlisa/generator"
	"log"
)

const InstagramClientId = "37b5984662894641bd7bbe4a597cd4b2"
const InstagramSecret = "9bd87bfb4408400d84c658d3ab5c728d"

func main() {
	instagramUsername := flag.String("username", "", "Instagram username for the user")
	inputFileName := flag.String("inputFile", "input.jpg", "Input file name")
	searchTerm := flag.String("searchTerm", "cats", "Stuff of which the mosaic is made of")

	flag.Parse()

	log.Printf("Starting out the mosaic generator with Instagram Username: %d, Input file: %d , Searc term %d ",
		*instagramUsername, *inputFileName, *searchTerm)

	generator.Generate(*instagramUsername, *searchTerm, *inputFileName)
}
