package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image/jpeg"
	"log"
	"net"
	"os"
)

type Message struct {
	ID   int
	Data []byte
}

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
	sendStruct(conn, imgbyte)
}

func sendStruct(conn net.Conn, imgByte []byte) {
	// lets create the message we want to send accross
	msg := Message{ID: 12, Data: imgByte}
	bin_buf := new(bytes.Buffer)

	// create a encoder object
	gobobj := gob.NewEncoder(bin_buf)
	// encode buffer and marshal it into a gob object
	gobobj.Encode(msg)

	fmt.Println(msg)

	conn.Write(bin_buf.Bytes())
}
