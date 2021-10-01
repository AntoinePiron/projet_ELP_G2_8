package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"time"
)

/**
Cette fonction va nous permettre d'analyser un bout de l'image à l'intérieur de la zone donnée par l'utilisateur
On indique également le fichier d'entrée et de sortie
*/
func analyze(upleftx int, uplefty int, width int, height int, input image.Image) {
	bounds := image.Rect(upleftx, uplefty, width, height)
	output := image.NewRGBA(bounds)

	for x := upleftx; x < width; x++ {
		for y := uplefty; y < height; y++ {
			oldPixel := input.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			output.Set(x, y, pixel)
		}
	}

	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, output, nil)
}

func main() {
	startTime := time.Now()
	file, err := os.Open("koala.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	analyze(0, 0, img.Bounds().Dx(), img.Bounds().Dy(), img)

	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())
}
