package main

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/Pixellettersfull-BnJ5.ttf
var fontData []byte

func loadFont(size float64) *text.GoTextFace {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: source,
		Size:   size,
	}
}

// Usage:
var MyFont = loadFont(40)
var myFontFace text.Face = MyFont
