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

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	r := &R{}
	fmt.Println(r)
	dec.Decode(r)
	fmt.Println(r)
	cmd := renderCommand(r)
	shout, err := net.Dial("tcp", ":8081")
	if err != nil {
		// handle shout error
	}
	shout.Write([]byte(cmd))
	fmt.Printf("redirect : %v\n", cmd)
}

func renderCommand(r *R) string {
	rDict := map[string]string{
		"houdini" : "hython",
	}
	runnable := rDict[r.Run]
	return fmt.Sprintf(`%s %s -c hou.node('%s').render()`, runnable, r.Scene, r.Driver)
}

