package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"time"
)

/**
Cette fonction va nous permettre d'analyser un bout de l'image à l'intérieur de la zone donnée par l'utilisateur
On indique également le fichier d'entrée et de sortie
*/
func analyze(upleftx int, uplefty int, width int, height int, input image.Image, final chan image.Image) {
	output := image.NewRGBA(input.Bounds())

	for x := upleftx; x < upleftx+width; x++ {
		for y := uplefty; y < uplefty+height; y++ {
			oldPixel := input.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			output.Set(x, y, pixel)
		}
	}
	final <- output
}

func main() {
	startTime := time.Now()
	file, err := os.Open("testFAT.JPG")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	finalImg := make(chan image.Image, 4)
	x := img.Bounds().Max.X / 2
	y := img.Bounds().Max.Y / 2
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			go analyze(x*i, y*j, img.Bounds().Dx()/2, img.Bounds().Dy()/2, img, finalImg)
		}
	}

	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	empty := image.NewRGBA(img.Bounds())
	capacity := cap(finalImg)
	for i := 0; i < capacity; i++ {
		outputIMG := <-finalImg
		draw.DrawMask(empty, empty.Bounds(), outputIMG, image.ZP, empty.Bounds(), image.ZP, draw.Over)
	}
	jpeg.Encode(outFile, empty, nil)

	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())
}
