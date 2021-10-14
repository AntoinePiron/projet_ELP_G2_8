package main

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
)

//On va se créer le go object qu'on va envoyer
type ImgForConn struct {
	id      int
	content image.Image
}

func main() {
	conn, errconn := net.Dial("tcp", "127.0.0.1:10000")
	defer conn.Close()
	if errconn != nil {
		os.Exit(1)
	} else {
		fmt.Println("Connection réussie")
	}

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

	sendStruct(conn, toSend)

}

func sendStruct(conn net.Conn, strct ImgForConn) {

	fmt.Println(strct)
	enc := gob.NewEncoder(conn)
	enc.Encode(strct)

}
