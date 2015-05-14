package lib

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func ResizeImages(context *AppContext, directoryName string) (int, bool) {
	dirDescriptor, err := os.Open("../downloads/" + directoryName)
	context.Log.Println("The directory name is :", directoryName)
	if err != nil {
		context.Log.Fatal("Unable to read directory.", err)
		return 0, false
	}
	defer dirDescriptor.Close()
	files, err := dirDescriptor.Readdir(-1)
	if err != nil {
		context.Log.Fatal("Unable to read files in the directrory")
		return 0, false
	}

	maxGoRoutines := 8
	fileNames := make(chan string, len(files))
	results := make(chan string, len(files))
	errorsChannel := make(chan error, 1)

	for _, fileObj := range files {
		if fileObj.Mode().IsRegular() {
			fileName := fileObj.Name()
			fileNames <- fileName
		}
	}

	close(fileNames)

	for j := 1; j <= maxGoRoutines; j++ {
		go extractAndProcess(context, fileNames, results, errorsChannel, directoryName)
	}

	for _, fileObj := range files {
		if fileObj.Mode().IsRegular() {
			select {
			case result := <-results:
				context.Log.Printf("The result for file: %v processing was %v \n", fileObj.Name(), result)
			case errMsg := <-errorsChannel:
				context.Log.Fatal("Sadly, something went wrong. Here's the error : ", errMsg)
			}
		}
	}

	return 0, true
}

func extractAndProcess(context *AppContext, fileNames <-chan string, results chan<- string, errChan chan<- error, directoryName string) {
	for fileName := range fileNames {
		imageFile, err := os.Open("../downloads/" + directoryName + "/" + fileName)
		if err != nil {
			context.Log.Printf("Unable to open the image file %v Error is %v \n", fileName, err)
			errChan <- err
		}
		img, _, err := image.Decode(imageFile)
		if err != nil {
			context.Log.Println("ERROR: Not able to decode the image file. Error is: ", err)
			errChan <- err
		}
		bounds := img.Bounds()

		boundsString := fmt.Sprintln(bounds)
		results <- boundsString
		imageFile.Close()
	}
}
