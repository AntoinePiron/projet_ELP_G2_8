package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, errconn := net.Dial("tcp", "127.0.0.1:10000")
	if errconn != nil {
		os.Exit(1)
	}
	fmt.Println("Connection réussie")
	time.Sleep(3 * time.Second)
	conn.Close()
	fmt.Println("Déco")
}
