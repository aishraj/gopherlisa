package main

import (
	"flag"
	"log"

	"github.com/aishraj/gopherlisa/generator"
	"github.com/aishraj/gopherlisa/server"
)

func main() {
	instagramUsername := flag.String("username", "", "Instagram username for the user")
	inputFileName := flag.String("inputFile", "input.jpg", "Input file name")
	searchTerm := flag.String("searchTerm", "cats", "Stuff of which the mosaic is made of")

	flag.Parse()

	log.Printf("Starting out the mosaic generator with Instagram Username: %s, Input file: %s , Search term %s ",
		*instagramUsername, *inputFileName, *searchTerm)
	go server.StartResponseServer()
	generator.Generate(*instagramUsername, *searchTerm, *inputFileName)
}
