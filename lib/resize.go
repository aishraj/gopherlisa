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
	//TODO: This is only placeholder stuff. Will add code later based on http://tech-algorithm.com/articles/nearest-neighbor-image-scaling/
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
	for _, fileObj := range files {
		if fileObj.Mode().IsRegular() {
			fileName := fileObj.Name()
			extractAndProcess(context, directoryName, fileName)
		}
	}
	return 0, true
}

func extractAndProcess(context *AppContext, directoryName, fileName string) {
	imageFile, err := os.Open("../downloads/" + directoryName + "/" + fileName)
	if err != nil {
		context.Log.Fatal("Unable to open the image file.", err)
	}
	defer imageFile.Close()
	img, _, err := image.Decode(imageFile)
	bounds := img.Bounds()

	fmt.Println(bounds)
	canvas := image.NewAlpha(bounds)
	op := canvas.Opaque()
	fmt.Println(op)
}
