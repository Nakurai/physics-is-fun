package main

import (
	"image/color"
	"math/rand"
)

func getRandomColor() color.RGBA {
	rColor := uint8(rand.Intn(MAX_COLOR-MIN_COLOR) + MIN_COLOR)
	gColor := uint8(rand.Intn(MAX_COLOR-MIN_COLOR) + MIN_COLOR)
	bColor := uint8(rand.Intn(MAX_COLOR-MIN_COLOR) + MIN_COLOR)

	return color.RGBA{
		R: rColor,
		G: gColor,
		B: bColor,
	}
}
