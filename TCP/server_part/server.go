package main

import (
	"fmt"
	"image"
	"net"
)

//On va se créer le go object qu'on va envoyer
type ImgForConn struct {
	id      int
	content image.Image
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

func handleConnection(connection net.Conn) {
	defer connection.Close()

}
