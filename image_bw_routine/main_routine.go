package main

import (
	"fmt"
	"image"
	"image/color"
	"time"
)

/**
Cette fonction va nous permettre d'analyser un bout de l'image à l'intérieur de la zone donnée par l'utilisateur
On indique également le fichier d'entrée et de sortie
*/
func analyze(upleftx int, uplefty int, width int, height int, input image.YCbCr, output image.RGBA) {
	for x := upleftx; x < width; x++ {
		for y := uplefty; y < height; y++ {
			oldPixel := input.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			output.Set(x, y, pixel)
		}
	}
}

func main() {
	startTime := time.Now()
	fmt.Println("Hello")
	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())
}
