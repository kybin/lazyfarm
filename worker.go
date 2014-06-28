package main

import (
	"net"
	"fmt"
)

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		msg := make([]byte, 2048)
		conn.Read(msg)
		fmt.Println(string(msg))
	}
}

