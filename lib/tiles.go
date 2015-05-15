package lib

import (
	"image"
	"image/color"
)

type Tile struct {
	X            int
	Y            int
	Rect         image.Rectangle
	AverageColor color.RGBA
	// SimilarImages used in 'matched' mosaic image
	MatchedImage string
}

func buildTilesLayer(targetWidth int, targetHeight int, tileSize int) [][]Tile {

	//TODO: Resize is needed before all this.
	// horzTiles, vertTiles := calcMosaicTiles(targetWidth, targetHeight, tileSize)
	// log.Printf("Tiles horizontal:%d vertical:%d", horzTiles, vertTiles)
	// // create a 2d array of imageTiles
	// imageTiles := make([][]Tile, horzTiles)
	// // Loop over the rows, allocating the slice for each row.
	// for i := range imageTiles {
	// 	imageTiles[i] = make([]Tile, vertTiles)
	// }
	//
	// // populate tiles with correct co-ordinates
	// for x := 0; x < horzTiles; x++ {
	// 	for y := 0; y < vertTiles; y++ {
	// 		currentTile := &imageTiles[x][y]
	// 		currentTile.X = x
	// 		currentTile.Y = y
	// 		tileStartX := x * tileSize
	// 		tileStartY := y * tileSize
	// 		tileEndX := tileStartX + tileSize
	// 		tileEndY := tileStartY + tileSize
	// 		// crop partial tile
	// 		if tileEndX >= targetWidth {
	// 			tileEndX = targetWidth
	// 		}
	// 		// crop partial tile
	// 		if tileEndY >= targetHeight {
	// 			tileEndY = targetHeight
	// 		}
	// 		currentTile.Rect = image.Rectangle{
	// 			image.Point{tileStartX, tileStartY},
	// 			image.Point{tileEndX, tileEndY},
	// 		}
	// 	}
	// }
	//
	// return imageTiles
	return nil
}
