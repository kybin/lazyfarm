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
	"flag"
	"strings"
)


func main() {
	var group string
	var server string
	flag.StringVar(&server, "server", "", "server address")
	flag.StringVar(&group, "group", "", "worker will serve this group of job")
	flag.Parse()
	if server == "" {
		fmt.Println("please specify server address")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var myaddr string = findMyAddress()
	go listenJob(myaddr, server, group)

	send(server, myaddr, "login", group)
	defer send(server, myaddr, "logout", group)

	go logoutAtExit(server, myaddr, group)

	for {
		time.Sleep(10*time.Second)
	}
}

func send(server, myaddr, status, group string) {
	conn, err := net.Dial("tcp", server)
	if err != nil{
		log.Fatal(err)
	}
	enc := gob.NewEncoder(conn)
	worker := &Worker{myaddr, group}
	err = enc.Encode(worker)
	if err != nil{
		log.Fatal(err)
	}
	err = enc.Encode(status)
	if err != nil{
		log.Fatal(err)
	}
}

func logoutAtExit(server, myaddr, group string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("interrupted...")
		send(server, myaddr, "logout", group)
		os.Exit(1)
	}()
}

func listenJob(myaddr, server, group string) {
	ln, err := net.Listen("tcp", myaddr)
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
		fmt.Printf("work start. (%v)\n", r)
		err = cmd.Run()
		if err != nil {
			log.Fatal(stderr.String())
		}
		fmt.Println(stdout.String())
		send(server, myaddr, "done", group)
		fmt.Println("work done.")
	}
}

func renderCommand(t *Task) *exec.Cmd {
	c := strings.Split(t.Cmd, " ")
	runnable := c[0]
	args := c[1:]
	return exec.Command(runnable, args...)
}

