/**
Code très largement inspiré du blog suivant :
https://mrwaggel.be/post/golang-transfer-a-file-over-a-tcp-socket/
*/

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()
	receivedFileName := receiveFileFromClient(connection)

	var wg sync.WaitGroup
	imgFile, err := os.Open(receivedFileName)
	if err != nil {
		fmt.Println("Probleme l57", err)
		return
	}
	defer imgFile.Close()
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		fmt.Println("Probleme l63", err)
		return
	}
	finalImg := image.NewRGBA(img.Bounds())
	const nbDiv = 8
	x := img.Bounds().Max.X / nbDiv
	y := img.Bounds().Max.Y / nbDiv
	for i := 0; i < nbDiv; i++ {
		for j := 0; j < nbDiv; j++ {
			//On oublie pas d'ajouter au waitGroup
			wg.Add(1)
			//On lance notre go routine
			go analyze(x*i, y*j, img.Bounds().Dx()/nbDiv, img.Bounds().Dy()/nbDiv, img, finalImg, &wg)
		}
	}
	wg.Wait()
	outFile, err := os.Create("changed.jpg")
	if err != nil {
		fmt.Println("Probleme l76", err)
		return
	}
	defer outFile.Close()
	jpeg.Encode(outFile, finalImg, nil)
	sendFileToClient(connection, "changed.jpg")
}

func sendFileToClient(conn net.Conn, name string) {
	//On ouvre le fichier et si jamais une erreur se produit on arrete la fonction avec le mot cle return
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("Probleme l87", err)
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

func receiveFileFromClient(connection net.Conn) string {
	bufferFileName := make([]byte, 64) //On fait correspondre les tailles avec l'envoie du cote server
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Probleme l126", err)
	}
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
	return fileName
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

//La fonction qui analyse l'image
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
