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
	"strconv"
	"flag"
)

var farmAddress string = ":8080"
var myAddress string = findMyAddress()

func findMyAddress() string {
	port := 8082
	for i :=0; i < 10; i++ {
		address := ":"+strconv.Itoa(port)
		ln, err := net.Listen("tcp", address)
		if err != nil {
			port++
			continue
		}
		ln.Close()
		fmt.Printf("my address is %v\n", address)
		return address
	}
	fmt.Printf("cannot find good port")
	os.Exit(1)
	return ""
}

func main() {
	var group string
	flag.StringVar(&group, "group", "", "worker will serve this group of job")
	flag.Parse()
	login()
	go logoutAtExit()
	go listenJob()
	for {
		time.Sleep(10*time.Second)
	}
}

func send(status string) {
	conn, err := net.Dial("tcp", farmAddress)
	if err != nil{
		log.Fatal(err)
	}
	enc := gob.NewEncoder(conn)
	worker := &Worker{myAddress}
	err = enc.Encode(worker)
	if err != nil{
		log.Fatal(err)
	}
	err = enc.Encode(status)
	if err != nil{
		log.Fatal(err)
	}
}


func login() {
	send("login")
}

func logout() {
	send("logout")
}

func logoutAtExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("interrupted...")
		logout()
		os.Exit(1)
	}()
}

func listenJob() {
	ln, err := net.Listen("tcp", myAddress)
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
		send("done")
		fmt.Println("work done.")
	}
}

func renderCommand(r *Task) *exec.Cmd {
	rDict := map[string]string{
		"houdini" : "hython",
	}
	runnable := rDict[r.Run]
	args := []string{r.Scene, "-c", fmt.Sprintf("hou.node('%s').render(frame_range=(%v,%v,1))", r.Driver, r.Frame, r.Frame)}
	return exec.Command(runnable, args...)
}

