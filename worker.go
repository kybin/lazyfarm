package main

import (
	"net"
	"log"
	"fmt"
	"encoding/gob"
	"bytes"
	"os/exec"
	"time"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	login()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("interrupted...")
		logout()
		os.Exit(1)
	}()
	for {
		time.Sleep(10*time.Second)
	}
}

func login() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil{
		log.Fatal(err)
	}
	enc := gob.NewEncoder(conn)
	worker := &Worker{":8081"}
	err = enc.Encode(worker)
	if err != nil{
		log.Fatal(err)
	}
	err = enc.Encode("login")
	if err != nil{
		log.Fatal(err)
	}
}

func logout() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil{
		log.Fatal(err)
	}
	enc := gob.NewEncoder(conn)
	worker := &Worker{":8081"}
	err = enc.Encode(worker)
	if err != nil {
		log.Fatal(err)
	}
	err = enc.Encode("logout")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logout")
}


func notuse() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		dec := gob.NewDecoder(conn)
		r := &Task{}
		dec.Decode(r)
		cmd := renderCommand(r)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(stderr.String())
		}
		fmt.Println(stdout.String())
		fmt.Println("work done.")
	}
}

func renderCommand(r *Task) *exec.Cmd {
	rDict := map[string]string{
		"houdini" : "hython",
	}
	runnable := rDict[r.Run]
	args := []string{r.Scene, "-c", fmt.Sprintf("hou.node('%s').render()", r.Driver)}
	return exec.Command(runnable, args...)
}

