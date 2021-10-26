/**
Code très largement inspiré du blog suivant :
https://mrwaggel.be/post/golang-transfer-a-file-over-a-tcp-socket/
*/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

//This constant can be anything from 1 to 65495, because the TCP package can only contain up to 65495 bytes of payload.
//It will define how big the chunks are of the file that we will send in bytes.
const BUFFERSIZE = 1024

//on définit le port comme une constante
const PORT = ":10000"

func main() {
	//On ouvre dans un premier temps le serveur TCP e vérifiant qu'il n'y a pas d'erreur
	server, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Problème lors de l'ouverture du serveur : ", err)
		os.Exit(1)
	}
	//On defer la fermeture du server
	defer server.Close()
	fmt.Println("Ouverture du serveur réussie")
	fmt.Println("En attente de connections ...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Problème lors de la connection du client : ", err)
		}
		fmt.Println("Client connecté")
		go sendFileToClient(connection)
	}
}

func sendFileToClient(conn net.Conn) {
	//On oublie pas de defer la fermture de la connection pour qu'elle se ferme automatiquement à la fin de l'execution
	defer conn.Close()
	//On ouvre le fichier et si jamais une erreur se produit on arrete la fonction avec le mot cle return
	file, err := os.Open("koala.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}

func receiveFileFromClient(connection net.Conn) {
	bufferFileName := make([]byte, 64) //On fait correspondre les tailles avec l'envoie du cote server
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
