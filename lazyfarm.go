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

func handleConnection(in net.Conn) {
	decoder := gob.NewDecoder(in)
	r := &R{}
	decoder.Decode(r)
	out, err := net.Dial("tcp", ":8081")
	encoder := gob.NewEncoder(out)
	err = encoder.Encode(r)
	if err != nil {
		// handle outconn error
	}
	fmt.Printf("redirect : %v\n", r)
}

