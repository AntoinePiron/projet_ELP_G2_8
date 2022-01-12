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

const BUFFERSIZE = 1024

var nbGoRoutine = 0

func main() {
	//Vérification de l'argument de l'utilisateur
	if len(os.Args) < 3 {
		fmt.Println("Veuillez deux arguments : port et nb de go routine")
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

	nbGoRoutine, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("nombre go routine incorrect")
		os.Exit(1)
	}
	if nbGoRoutine < 1 {
		fmt.Println("Veuillez rentrer une valeur positive de go routines")
		os.Exit(1)
	}

	//On ouvre dans un premier temps le serveur TCP e vérifiant qu'il n'y a pas d'erreur
	server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Problème lors de l'ouverture du serveur : ", err)
		os.Exit(1)
	}
	//On defer la fermeture du server = on le ferme une fois le main terminé
	defer server.Close()
	fmt.Println("Ouverture du serveur réussie")
	fmt.Println("En attente de connections ...")
	numberOfConnections := 0
	//On créé un boucle infinie attendant les différents connections
	for {
		//On recpèère une connection et on gere une erreur si nécessaire
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Problème lors de la connection du client : ", err)
			continue
		}
		fmt.Println("Client connecté")
		numberOfConnections += 1
		//Quand la connection à réussi on va la traiter avec une go routine
		go handleConnection(connection, numberOfConnections)
	}
}

/**
La fonction qui permet de gérer chaque utilisateur
paramètre :
 - connection --> la connection de l'utilisateur de type net.Conn
*/
func handleConnection(connection net.Conn, numberOfConnections int) {
	outName := "changed_" + strconv.Itoa(numberOfConnections) + ".jpg"
	//On ferme la connection une fois toutes la méthode finie
	defer connection.Close()
	//On traite le fichier et on récupère son nom
	receivedFileName := receiveFileFromClient(connection)
	//On s'occuper alors de modifier notre image
	imageProcess(receivedFileName, outName)
	//une fois finit on renvoie la nouvelle image au client
	sendFileToClient(connection, outName)
}

/**
La fonction qui permet d'envoyer un fichier au client
paramètres :
 - conn --> la connection du client du type net.Conn
 - name --> le nom du fichier à envoyer du type string
*/
func sendFileToClient(conn net.Conn, name string) {
	//On ouvre le fichier et si jamais une erreur se produit on arrete la fonction avec le mot cle return
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	//On recupère les statistiques du fichier
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	//On récupère la taille et le nom du fichier pour les envoyer directement à l'utilisateur
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	//Un fois réalisé on peut alors envoyer directement notre image bout par bout car on est limité à une certaine taille de buffer avec une connection TCP
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

/**
La fonction qui permet de recevoir un fichier du client
paramètre :
 - connection --> la connection du client du type net.Conn
return :
 - fileName --> le nom du fichier reçu
*/
func receiveFileFromClient(connection net.Conn) string {
	//On recoit dans un premier temps le nom et la taille du fichier selon les tailles prédéfini (On remarquera qu'elles coresspondent également à notre envoie)
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	//On peut alors créerle fichier de sortie
	newFile, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
	}
	/**Initialize a 64 bit integer that will keep count of how many bytes we have
	received so far, so we can tell when to stop reading the chunks from the server.*/
	var receivedBytes int64
	//On remplit alors ce fameux fichier avec les données reçu
	for {
		if (fileSize - receivedBytes) <= BUFFERSIZE {
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

/**
Cette fonction est un peu particulière
Elle permet de remplir un paquet avec la taille voulue si il n'est pas complet
cad que si notre string ne fait que 10 bits alors qu'on veut emvoyer un paquet de 64bits cette fonction va se charger de le remplir avec le carac ":" pour atteindre la taille souhaité
paramètre :
 - returnString --> notre string de base qu'on veut "épaissir"
 - toLength --> la longueur souhaitée
*/
func fillString(returnString string, toLength int) string {
	for {
		lengtString := len(returnString)
		if lengtString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}

/**
La fonction qui permet de déclencher le traitement de l'image en la sous-traitant avec plusieurs go routines
elle est détallé dans le dossier /image_bw_routine pour la comprendre
paramètre :
 - imageName --> le nom de l'image reçu de l'utilisateur
*/
func imageProcess(imageName string, outName string) {
	var wg sync.WaitGroup
	imgFile, err := os.Open(imageName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer imgFile.Close()
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	finalImg := image.NewRGBA(img.Bounds())
	x := img.Bounds().Max.X / nbGoRoutine
	y := img.Bounds().Max.Y
	for i := 0; i < nbGoRoutine; i++ {
		//On oublie pas d'ajouter au waitGroup
		wg.Add(1)
		//On lance notre go routine
		go analyze(x*i, 0, img.Bounds().Dx()/nbGoRoutine, y, img, finalImg, &wg)
	}
	wg.Wait()
	outFile, err := os.Create(outName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	jpeg.Encode(outFile, finalImg, nil)
}

/**
La fonction qui analyse une portion de l'image
paramètres :
 - upleftx --> La coordonné x du coin haut gauche de la zone à traiter
 - uplefty --> La coordonné y du coin haut gauche de la zone à traiter
 - width --> la largeur de la zone à traiter
 - height --> la hauteur de la zone à traiter
 - input --> l'image de base
 - final --> l'image de sortie ou on va écrire les nouveaux pixels
 - wg --> le waitgroup permettant de gérer les go routines
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
