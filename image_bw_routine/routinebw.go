package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"strconv"
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
	//Vérification de l'argument de l'utilisateur
	if len(os.Args) < 3 {
		fmt.Println("Veuillez deux arguments : nombres de divisions et nom de fichier")
		os.Exit(1)
	}
	nbDiv, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Input incorrecte")
		os.Exit(1)
	}
	if nbDiv <= 0 {
		fmt.Println("Veuillez rentrer une valeur positive de division")
		os.Exit(1)
	}
	filename := os.Args[2]
	var wg sync.WaitGroup //On initialise notre waitgroup pour notre travail de goroutine par la suite

	//Ce premier bloc permet d'ouvrir notre image sous forme de file et de vérifier au'il n'y a aucune erreur
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Depuis notre fichier on le convertit en image pour go, on vérifie à nouveau qu'il n'y a pas d'erreur
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	//On vient créer une image vide faisant la même taille que l'image d'origine
	finalImg := image.NewRGBA(img.Bounds())
	//il s'agit ici de la longueur d'une sous section
	x := img.Bounds().Max.X / nbDiv
	//Il s'agit ici de la largeur d'une sous section
	y := img.Bounds().Max.Y / nbDiv
	//On va sauvegarder le tps avant l'execution des go routines pour voir le temps d'execution du programme
	startTime := time.Now()
	//Le double for va permettre de lancer les différentes go routines pour les sous segments
	for i := 0; i < nbDiv; i++ {
		for j := 0; j < nbDiv; j++ {
			//On oublie pas d'ajouter au waitGroup
			wg.Add(1)
			//On lance notre go routine
			go analyze(x*i, y*j, img.Bounds().Dx()/nbDiv, img.Bounds().Dy()/nbDiv, img, finalImg, &wg)
		}
	}
	//Le wait va permettre d'eviter l'avancée du programme tant que toutes les go routines ne se sont pas executé
	//Cela permet d'être sur que l'image est bien constitué en entier et au'on peut l'exporter
	wg.Wait()

	//On releve par ailleurs le temps total d'execution qu'on affiche
	//On ne prend pas en compte le temps I/O qui ne dépend pas des go routines
	totalTime := time.Since(startTime)
	fmt.Println("Durée totale : " + totalTime.String())

	//Finalement on genere notre fichier de sortie
	outFile, err := os.Create("changed.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	//Et on le sort dans le bon format
	jpeg.Encode(outFile, finalImg, nil)
}
