package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net"
	"os"
)

func main() {
	conn, errconn := net.Dial("tcp", "127.0.0.1:10000")
	defer conn.Close()
	if errconn != nil {
		os.Exit(1)
	} else {
		fmt.Println("Connection réussie")
	}

	//Ce premier bloc permet d'ouvrir notre image sous forme de file et de vérifier au'il n'y a aucune erreur
	file, err := os.Open("koala.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Depuis notre fichier on le convertit en image pour go, on vérifie à nouveau qu'il n'y a pas d'erreur
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}
	imgbyte := buf.Bytes()
	conn.Write(imgbyte)
}
