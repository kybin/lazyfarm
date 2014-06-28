package main

import (
	"net"
	"fmt"
)

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
		buf := make([]byte, 2048)
		conn.Read(buf)
		shout, err := net.Dial("tcp", ":8081")
		if err != nil {
			// handle shout error
		}
		shout.Write(buf)
		fmt.Printf("redirect : %v\n", string(buf))
	}
}
