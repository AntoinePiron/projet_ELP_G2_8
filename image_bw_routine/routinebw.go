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
func analyze(upleftx int, uplefty int, width int, height int, input image.Image, final *image.RGBA, wg *sync.WaitGroup) {
	//Le double for permet de se déplacer parmi tt les pixels de la zone
	for x := upleftx; x < upleftx+width; x++ {
		for y := uplefty; y < uplefty+height; y++ {
			//Le pixel présent à l'origine sur l'image
			oldPixel := input.At(x, y)
			//On prend ses valeurs RGB, on ignore l'alpha non utile
			r, g, b, _ := oldPixel.RGBA()
			//Cette conversion permet de passer d'une valeur RGB à un niveau de gris
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			//On genere alors notrpe pixel
			pixel := color.Gray{uint8(lum / 256)}
			//Puis on écrit le pixel à la même coordonne dans le fichier de sortie
			final.Set(x, y, pixel)
		}
	}
	//Une fois le travail execute on le renseigne au WaitGroup
	wg.Done()
}

func main() {
	var wg sync.WaitGroup //On initialise notre waitgroup pour notre travail de goroutine par la suite

	//Ce premier bloc permet d'ouvrir notre image sous forme de file et de vérifier au'il n'y a aucune erreur
	file, err := os.Open("testFAT.JPG")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Depuis notre fichier on le convertit en image pour go, on vérifie à nouveau qu'il n'y a pas d'erreur
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	//On vient créer une image vide faisant la même taille que l'image d'orugine
	finalImg := image.NewRGBA(img.Bounds())
	const nbDiv = 8
	x := img.Bounds().Max.X / nbDiv
	y := img.Bounds().Max.Y / nbDiv
	startTime := time.Now()
	for i := 0; i < nbDiv; i++ {
		for j := 0; j < nbDiv; j++ {
			wg.Add(1)
			go analyze(x*i, y*j, img.Bounds().Dx()/nbDiv, img.Bounds().Dy()/nbDiv, img, finalImg, &wg)
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
