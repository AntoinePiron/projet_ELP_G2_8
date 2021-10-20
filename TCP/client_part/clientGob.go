package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Message struct {
	ID   int
	Data string
}

func main() {
	conn, errconn := net.Dial("tcp", "127.0.0.1:10000")
	defer conn.Close()
	if errconn != nil {
		os.Exit(1)
	} else {
		fmt.Println("Connection réussie")
	}

	sendStruct(conn)
}

func sendStruct(conn net.Conn) {
	// lets create the message we want to send accross
	msg := Message{ID: 12, Data: "coucou"}
	bin_buf := new(bytes.Buffer)

	// create a encoder object
	gobobj := gob.NewEncoder(bin_buf)
	// encode buffer and marshal it into a gob object
	gobobj.Encode(msg)

	fmt.Println(msg)

	conn.Write(bin_buf.Bytes())
}
