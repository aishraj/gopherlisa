package lib

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

func CreateMosaic(context *AppContext, srcName, destDirName string) image.Image {
	srcImg, err := LoadImage("/tmp/" + srcName + ".jpg")
	if err != nil {
		context.Log.Println("Unable to open the input file. Error is", err)
		return nil
	}
	sourceImage := ToNRGBA(srcImg)
	outputImageWidth := 600
	outputImageHeight := calcRelativeImageHeight(sourceImage.Bounds().Max.X, sourceImage.Bounds().Max.Y, outputImageWidth)

	resizedImage := Resize(sourceImage, outputImageWidth, outputImageHeight)
	// how many tiles?
	imageTiles := createTiles(outputImageWidth, outputImageHeight)
	// analyse input image colours
	analysedTiles := analyseImageTileColours(resizedImage, imageTiles)

	// update tiles with details of similar images
	preparedTiles := updateSimilarColourImages(context, analysedTiles, srcName)

	// draw photo tiles
	photoImage := drawPhotoTiles(resizedImage, &preparedTiles, 64)

	outputImagePath := "/tmp/output.jpeg"
	// save image created
	err = SaveImage(outputImagePath, &photoImage)
	return nil
}
func SaveImage(imagePath string, imageToSave *image.Image) error {
	if imgFilePng, err := os.Create(imagePath); err != nil {
		log.Printf("Error saving PNG image: %s\n", err)
		return err
	} else {
		defer imgFilePng.Close()
		buffer := bufio.NewWriter(imgFilePng)
		var opt jpeg.Options

		opt.Quality = 80
		err := jpeg.Encode(buffer, *imageToSave, &opt)
		if err != nil {
			log.Printf("Error encoding image:%s", err)
			return err
		}
		buffer.Flush()

		return nil
	}
}
func LoadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Error during LoadImage: %s", err)
		return nil, err
	}
	defer file.Close()
	loadedImage, _, err := image.Decode(file)

	return loadedImage, err

}

func calcRelativeImageHeight(originalWidth int, originalHeight int, targetWidth int) int {
	floatWidth := float64(originalWidth)
	floatHeight := float64(originalHeight)

	aspectRatio := float64(targetWidth) / floatWidth

	adjustedHeight := floatHeight * aspectRatio

	targetHeight := int(adjustedHeight)
	log.Printf("Source width:%d height:%d Target width:%d height:%d\n", originalWidth, originalHeight, targetWidth, targetHeight)
	return targetHeight
}

func createTiles(targetWidth int, targetHeight int) [][]Tile {

	tileSize := 64

	horzTiles := targetWidth / tileSize
	if targetWidth%tileSize > 0 {
		horzTiles++
	}
	vertTiles := targetHeight / tileSize
	if targetHeight%tileSize > 0 {
		vertTiles++
	}

	log.Printf("Tiles horizontal:%d vertical:%d", horzTiles, vertTiles)
	// create a 2d array of imageTiles
	imageTiles := make([][]Tile, horzTiles)
	// Loop over the rows, allocating the slice for each row.
	for i := range imageTiles {
		imageTiles[i] = make([]Tile, vertTiles)
	}

	// populate tiles with correct co-ordinates
	for x := 0; x < horzTiles; x++ {
		for y := 0; y < vertTiles; y++ {
			currentTile := &imageTiles[x][y]
			currentTile.X = x
			currentTile.Y = y
			tileStartX := x * tileSize
			tileStartY := y * tileSize
			tileEndX := tileStartX + tileSize
			tileEndY := tileStartY + tileSize
			// crop partial tile
			if tileEndX >= targetWidth {
				tileEndX = targetWidth
			}
			// crop partial tile
			if tileEndY >= targetHeight {
				tileEndY = targetHeight
			}
			currentTile.Rect = image.Rectangle{
				image.Point{tileStartX, tileStartY},
				image.Point{tileEndX, tileEndY},
			}
		}
	}

	return imageTiles
}

func analyseImageTileColours(sourceImage image.Image, imageTiles [][]Tile) [][]Tile {
	for _, tiles := range imageTiles {
		for _, tile := range tiles {
			tile.AverageColor = findAverageColor(sourceImage, tile.Rect)
			imageTiles[tile.X][tile.Y].AverageColor = tile.AverageColor
		}
	}

	return imageTiles
}

func findAverageColor(sourceImage image.Image, targetRect image.Rectangle) color.RGBA {
	croppedImage := Crop(sourceImage, targetRect)
	return FindProminentColour(croppedImage)

}
func updateSimilarColourImages(context *AppContext, imageTiles [][]Tile, indexName string) [][]Tile {

	for _, tiles := range imageTiles {
		for _, tile := range tiles {

			imageTiles[tile.X][tile.Y].MatchedImage = findClosestMatch(context, tile.AverageColor, indexName)

		}
	}

	return imageTiles
}

func drawPhotoTiles(sourceImage image.Image, imageTiles *[][]Tile, tileWidth int) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	photoImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(photoImage, photoImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	for _, tiles := range *imageTiles {
		for _, tile := range tiles {

			// draw image using first tile discovered
			if tile.MatchedImage != "" {
				//

				tileImage, err := LoadImage(tile.MatchedImage)
				if err != nil {
					panic("Error loading image")
				}
				tileImageNRGBA := ToNRGBA(tileImage)
				// resize image to tile size
				resizedImage := Resize(tileImageNRGBA, tileWidth, tileWidth)
				draw.Draw(photoImage, tile.Rect, resizedImage, tileImage.Bounds().Min, draw.Src)

			}

		}
	}

	return photoImage
}
