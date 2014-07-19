package main

import (
	"net"
	"fmt"
	"encoding/gob"
)

type R struct {
	Run, Scene, Driver, Frames string
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(inconn net.Conn) {
	dec := gob.NewDecoder(inconn)
	r := &R{}
	dec.Decode(r)
	outconn, err := net.Dial("tcp", ":8081")
	enc := gob.NewEncoder(outconn)
	err = enc.Encode(r)
	if err != nil {
		// handle outconn error
	}
	fmt.Printf("redirect : %v\n", r)
}

