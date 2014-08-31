package main

import (
	"time"
	"strings"
	"fmt"
	"log"
	"net"
	"encoding/gob"
	"errors"
)

func main() {
	msgchan := make(chan string)
	popchan := make(chan string)

	go workerStack(msgchan, popchan)

	go listenWorker(msgchan)
	go listenJob(msgchan, popchan)

	for {
		time.Sleep(time.Second)
		fmt.Println("tired...")
	}
}

func workerStack(msgchan chan string, popchan chan string) {
	stack := make([]string, 0)
	var msg string
	var msgs []string
	var status string
	var address string
	for {
		msg = <-msgchan
		msgs = strings.Split(msg, " ")
		status = msgs[0]
		switch status {
		case "login": // push
			address = msgs[1]
			stack = append(stack, address)
		case "logout": // find index and delete
			address = msgs[1]
			idx := -1
			for i, v := range stack {
				if (v == address) {
					idx = i
					break
				}
			}
			if idx == -1 {
				notfound := errors.New("not found " + address)
				log.Fatal(notfound)
			}
			stack = append(stack[:idx], stack[idx+1:]...) // delete address from stack
		case "waiting": // same with login yet.
			address = msgs[1]
			stack = append(stack, address)
		case "need": // pop
			if len(stack) == 0 {
				address = ""
			} else {
				last := len(stack)-1
				address = stack[last]
				stack = stack[:last]
			}
			popchan <- address
		default:
			notexpect := errors.New(fmt.Sprintf("not expected status '%v'", status))
			log.Fatal(notexpect)
		}
		fmt.Printf("%v\n", stack)
	}
}


func listenWorker(msgchan chan string) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	// 작업자 상태 변경 플롯
	// 작업자 추가 - login
	// 작업자 삭제 - logout
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connected")
		worker, status := handleWorkerConn(conn)
		switch status {
		case "login", "logout":
			msgchan <- status + " " + worker.Address
		default:
			log.Fatal("unknown status")
		}
	}
}

func handleWorkerConn(conn net.Conn) (*Worker, string)  {
	decoder := gob.NewDecoder(conn)
	w := &Worker{}
	err := decoder.Decode(w)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("worker decoded")
	var status string
	err = decoder.Decode(&status)
	if err != nil {
		log.Fatal(err)
	}
	return w, status
}


func listenJob(msgchan chan string, popchan chan string) {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		job := &Job{}
		decoder := gob.NewDecoder(conn)
		decoder.Decode(job)
		go handleJob(job, msgchan, popchan)
	}
}

func handleJob(job *Job, msgchan chan string, popchan chan string) {
	tasks := jobToTasks(job)
	for t := range tasks {
		msgchan <- "need"
		worker_address := <-popchan
		out, err := net.Dial("tcp", worker_address)
		encoder := gob.NewEncoder(out)
		err = encoder.Encode(t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("redirect : %v\n", t)
	}
}

func jobToTasks(job *Job) []Task {
	tasks := make([]Task, 0)
	return tasks
}
