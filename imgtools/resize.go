package imgtools

import (
	"fmt"
	"github.com/aishraj/gopherlisa/common"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

func ResizeImages(context *common.AppContext, directoryName string) (int, bool) {
	dirDescriptor, err := os.Open(common.DownloadBasePath + directoryName) //TODO change this
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

	successCount := 0

	for _, fileObj := range files {
		if fileObj.Mode().IsRegular() {
			select {
			case result := <-results:
				successCount++
				context.Log.Printf("The result for file: %v processing was %v \n", fileObj.Name(), result)
			case errMsg := <-errorsChannel:
				context.Log.Fatal("Sadly, something went wrong. Here's the error : ", errMsg)
			}
		}
	}

	return successCount, true
}

func extractAndProcess(context *common.AppContext, fileNames <-chan string, results chan<- string, errChan chan<- error, directoryName string) {
	for fileName := range fileNames {
		filePath := common.DownloadBasePath + directoryName + "/" + fileName
		imageFile, err := os.Open(filePath)
		if err != nil {
			context.Log.Printf("Unable to open the image file %v Error is %v \n", fileName, err)
			errChan <- err
		}
		img, _, err := image.Decode(imageFile)
		if err != nil {
			context.Log.Println("ERROR: Not able to decode the image file. Error is: ", err)
			errChan <- err
		}

		nrgbaImage := ToNRGBA(img)
		bounds := nrgbaImage.Bounds()
		boundsString := fmt.Sprintln(bounds)
		context.Log.Println("The bounds BEFORE the resize are: ", boundsString)

		nrgbaImage = Resize(nrgbaImage, 64, 64)

		bounds = nrgbaImage.Bounds()
		boundsString = fmt.Sprintln(bounds)
		context.Log.Println("The bounds AFTER the resize are: ", boundsString)

		var opt jpeg.Options

		opt.Quality = 80

		imageFile.Close()
		writeFile, err := os.Create(filePath)
		if err != nil {
			context.Log.Println("Unable to open the file for writing. Error is:", err)
		}
		err = jpeg.Encode(writeFile, nrgbaImage, &opt)
		if err != nil {
			context.Log.Println("ERROR: Not able to write to file with JPEG encoding", err)
			errChan <- err
		}

		results <- boundsString
		writeFile.Close()
	}
}

func ToNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return CloneImage(img)
}

func Resize(src *image.NRGBA, width, height int) *image.NRGBA {
	dstW, dstH := width, height

	srcBounds := src.Bounds()
	srcW := srcBounds.Max.X
	srcH := srcBounds.Max.Y

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))

	dx := float64(srcW) / float64(dstW)
	dy := float64(srcH) / float64(dstH)

	partStart := 0
	partEnd := dstH

	for dstY := partStart; dstY < partEnd; dstY++ {
		fy := (float64(dstY)+0.5)*dy - 0.5

		for dstX := 0; dstX < dstW; dstX++ {
			fx := (float64(dstX)+0.5)*dx - 0.5

			srcX := int(math.Min(math.Max(math.Floor(fx+0.5), 0.0), float64(srcW)))
			srcY := int(math.Min(math.Max(math.Floor(fy+0.5), 0.0), float64(srcH)))

			srcOff := srcY*src.Stride + srcX*4
			dstOff := dstY*dst.Stride + dstX*4

			copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
		}
	}

	return dst
}

func CloneImage(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	dstBounds := srcBounds.Sub(srcBounds.Min)
	dstW := dstBounds.Dx()
	dstH := dstBounds.Dy()
	dst := image.NewNRGBA(dstBounds)

	partStart := 0
	partEnd := dstH

	for dstY := partStart; dstY < partEnd; dstY++ {
		di := dst.PixOffset(0, dstY)
		for dstX := 0; dstX < dstW; dstX++ {

			c := color.NRGBAModel.Convert(img.At(srcMinX+dstX, srcMinY+dstY)).(color.NRGBA)
			dst.Pix[di+0] = c.R
			dst.Pix[di+1] = c.G
			dst.Pix[di+2] = c.B
			dst.Pix[di+3] = c.A

			di += 4

		}
	}

	return dst
}
