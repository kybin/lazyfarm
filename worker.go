package main

import (
	"net"
	"log"
	"fmt"
	"strings"
	"bytes"
	"os/exec"
)

func main() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		msg := make([]byte, 2048)
		n, err := conn.Read(msg)
		if err != nil {
			// handle error
		}
		cmdlist := strings.Split(string(msg[:n]), " ")
		name := cmdlist[0]
		args := cmdlist[1:]
		fmt.Println(name)
		fmt.Println(args)
		cmd := exec.Command(name, args...)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(stderr.String())
		}
		fmt.Println(stdout.String())
	}
}

