package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"
)

/**
Cette fonction va nous permettre d'analyser un bout de l'image à l'intérieur de la zone donnée par l'utilisateur
On indique également le fichier d'entrée et de sortie
*/
func analyze(upleftx int, uplefty int, width int, height int, input image.Image, final *image.RGBA, wg *sync.WaitGroup, m *sync.Mutex) {

	for x := upleftx; x < upleftx+width; x++ {
		for y := uplefty; y < uplefty+height; y++ {
			oldPixel := input.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			final.Set(x, y, pixel)
		}
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex
	file, err := os.Open("testFAT.JPG")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	finalImg := image.NewRGBA(img.Bounds())
	const nbDiv = 4
	x := img.Bounds().Max.X / nbDiv
	y := img.Bounds().Max.Y / nbDiv
	startTime := time.Now()
	for i := 0; i < nbDiv; i++ {
		for j := 0; j < nbDiv; j++ {
			wg.Add(1)
			go analyze(x*i, y*j, img.Bounds().Dx()/nbDiv, img.Bounds().Dy()/nbDiv, img, finalImg, &wg, &m)
		}
	}
	wg.Wait()
	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())

	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, finalImg, nil)
}
