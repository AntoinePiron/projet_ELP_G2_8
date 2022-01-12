/**
Code très largement inspiré du blog suivant :
https://mrwaggel.be/post/golang-transfer-a-file-over-a-tcp-socket/
La partie vlient utilise globalement les mêmes méthodes que la partie serveur, les méthodes sont donc détaillé dans
la partie serveur
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

const BUFFERSIZE = 1024

func main() {
	sendFileName := ""
	if len(os.Args) < 3 {
		fmt.Println("Veuillez deux arguments : port et nom de fichier")
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Port incorrect")
		os.Exit(1)
	}
	if port <= 1024 {
		fmt.Println("Veuillez rentrer une valeur compatible de port (>1024)")
		os.Exit(1)
	}
	sendFileName = os.Args[2]
	connection, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Println("Connected to server")

	sendFileToServer(connection, sendFileName)
	receiveFileFromServer(connection)
}

func sendFileToServer(conn net.Conn, name string) {
	//On ouvre le fichier et si jamais une erreur se produit on arrete la fonction avec le mot cle return
	file, err := os.Open(name)
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
	fmt.Println("File has been sent")
	return
}

func receiveFileFromServer(connection net.Conn) {
	bufferFileName := make([]byte, 64) //On fait correspondre les tailles avec l'envoie du cote server
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
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
