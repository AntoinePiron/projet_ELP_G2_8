package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
	"time"
)

//On va se créer le go object qu'on va envoyer
type ImgForConn struct {
	id      int
	content image.Image
}

func main() {
	file, err := os.Open("koala.jpg")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	toSend := ImgForConn{
		0,
		img,
	}
	fmt.Println(toSend)

	conn, errconn := net.Dial("tcp", "127.0.0.1:10000")
	if errconn != nil {
		os.Exit(1)
	} else {
		fmt.Println("Connection réussie")
	}
	time.Sleep(3 * time.Second)
	conn.Close()
	fmt.Println("Déco")
}
