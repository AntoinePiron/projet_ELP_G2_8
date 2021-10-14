package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type Message struct {
	ID   string
	Data string
}

func main() {
	const port = ":10000"
	ln, err := net.Listen("tcp", port) //avec port le numéro du port, ln = listener
	if err != nil {
		fmt.Print("Problème lors de l'ouverture du serveur")
		panic(err)
	} else {
		fmt.Println("Ouverture du serveur réussie")
	}
	for {
		fmt.Println("En attente de connection ...")
		conn, errconn := ln.Accept() //On accepte la connection et on met l'identifiant de la session dans conn
		//Cette ligne bloque le code tant qu'il n'y a pas de connectiom
		if errconn != nil {
			panic(errconn)
		} else {
			fmt.Println("Connection réussie")
		}
		//On prend tout de suite en charge la connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// create a temp buffer
	tmp := make([]byte, 500)

	// loop through the connection to read incoming connections. If you're doing by
	// directional, you might want to make this into a seperate go routine
	for {
		_, err := conn.Read(tmp)
		if err != nil {
			break
		}

		// convert bytes into Buffer (which implements io.Reader/io.Writer)
		tmpbuff := bytes.NewBuffer(tmp)
		tmpstruct := new(Message)

		// creates a decoder object
		gobobj := gob.NewDecoder(tmpbuff)
		// decodes buffer and unmarshals it into a Message struct
		gobobj.Decode(tmpstruct)

		// lets print out!
		fmt.Println(tmpstruct)
		return
	}

}
