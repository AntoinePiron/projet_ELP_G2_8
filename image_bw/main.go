package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

type ImageSet interface {
	Set(x, y int, c color.Color)
}

func main() {
	file, err := os.Open("koala.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()

	imgSet := image.NewRGBA(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			oldPixel := img.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			imgSet.Set(x, y, pixel)
		}
	}

	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, imgSet, nil)

}
