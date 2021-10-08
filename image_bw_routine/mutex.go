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
			//output.Set(x, y, pixel)
			m.Lock()
			final.Set(x, y, pixel)
			m.Unlock()
		}
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex

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

	finalImg := image.NewRGBA(img.Bounds())
	x := img.Bounds().Max.X / 2
	y := img.Bounds().Max.Y / 2
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go analyze(x*i, y*j, img.Bounds().Dx()/2, img.Bounds().Dy()/2, img, finalImg, &wg, &m)
		}
	}
	wg.Wait()

	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, finalImg, nil)

	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())
}
