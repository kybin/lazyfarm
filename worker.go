package main

import (
	"net"
	"log"
	"fmt"
	"encoding/gob"
	// "bytes"
	"os/exec"
	"time"
	"os"
	"os/signal"
	"syscall"
	"flag"
	"strings"
	"bufio"
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

	// make a log dir
	os.Mkdir("log", 0755)

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
			log.Fatal(err)
		}
		dec := gob.NewDecoder(conn)
		r := &Task{}
		dec.Decode(r)
		cmd := renderCommand(r)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("work start. (%v)\n", r)
		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.OpenFile("log/testlog.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600) 
		if err != nil {
			panic(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			_, err = f.WriteString(scanner.Text()+"\n")
			if err != nil {
				panic(err)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		// fmt.Println(stdout.String())
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

