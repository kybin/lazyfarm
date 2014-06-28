package main

import (
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		// handle error
	}
	msg := os.Args[1]
	conn.Write([]byte(msg))
}
