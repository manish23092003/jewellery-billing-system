package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func generateTextWatermark(text string) []byte {
	if text == "" {
		text = "CONFIDENTIAL"
	}
	
	// A4 aspect ratio
	width := 800
	height := 1131
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	col := color.RGBA{R: 210, G: 210, B: 210, A: 100}
	
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
	}
	
	textWidth := len(text) * 7
	stepX := textWidth + 80
	stepY := 60
	
	for y := 20; y < height+stepY; y += stepY {
		offsetX := 0
		if (y/stepY)%2 != 0 {
			offsetX = stepX / 2
		}
		
		for x := -offsetX; x < width+stepX; x += stepX {
			d.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
			d.DrawString(text)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func main() {
	b := generateTextWatermark("AWESOME JEWELLERS")
	os.WriteFile("test_watermark_grid.png", b, 0644)
	fmt.Printf("Generated %d bytes\n", len(b))
}
