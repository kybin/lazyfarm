package main

import (
	"time"
	//"strings"
	"fmt"
	"log"
	"net"
	"encoding/gob"
	"errors"
)

var workerStackChan = make(chan WorkerStackMsg)

func main() {
	go workerStack()
	go listenWorker()
	go listenJob()

	for {
		time.Sleep(time.Second)
		// fmt.Println("tired...")
	}
}

func workerStack() {
	stack := make([]string, 0)
	for {
		msg := <-workerStackChan

		switch msg.Type {
		case "push":
			stack = append(stack, msg.WorkerAddress)
		case "pop":
			address := ""
			if len(stack) != 0 {
				last := len(stack)-1
				address = stack[last]
				stack = stack[:last]
			}
			msg.Reply <- address
		case "delete":
			di := -1
			for i, v := range stack {
				if (v == msg.WorkerAddress) {
					di = i
					break
				}
			}
			if di == -1 {
				log.Fatal(errors.New("worker not found " + msg.WorkerAddress))
			}
			stack = append(stack[:di], stack[di+1:]...)
		default:
			log.Fatal(errors.New(fmt.Sprintf("not expected message type '%v'", msg.Type)))
		}
		fmt.Printf("worker list - %v\n", stack)
	}
}


func listenWorker() {
	ip, err := localIP()
	if err != nil {
		log.Fatal(err)
	}

	// worker socket
	ln, err := net.Listen("tcp", ip.String()+":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening worker message from", ip.String()+":8080")

	for {
		// data in from worker socket
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// decode worker's status and infomation
		decoder := gob.NewDecoder(conn)
		worker := Worker{}
		err = decoder.Decode(&worker)
		if err != nil {
			log.Fatal(err)
		}
		var status string
		err = decoder.Decode(&status)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("worker %v - %v\n", status, worker)

		handleWorker(status, worker)
	}
}

func handleWorker(status string, worker Worker) {
	var msgtype string

	switch status {
	case "login":
		msgtype = "push"
	case "logout":
		msgtype = "delete"
	case "done":
		msgtype = "push"
	default:
		log.Fatal("unknown status")
	}

	workerStackChan <- WorkerStackMsg{Type:msgtype, WorkerAddress:worker.Address}
}

func listenJob() {
	ip, err := localIP()
	if err != nil {
		log.Fatal(err)
	}

	// job socket
	ln, err := net.Listen("tcp", ip.String()+":8081")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening job message from", ip.String()+":8081")

	for {
		// data in from job socket
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// job decoding
		job := &Job{}
		decoder := gob.NewDecoder(conn)
		decoder.Decode(job)
		fmt.Printf("job added - %v\n", job)

		go handleJob(job)
	}
}

func handleJob(job *Job) {
	// separate it to tasks
	tasks, err := jobToTasks(job)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tasks {
		fmt.Printf("%v\n", t)
		worker_address := findWorker()
		sendTask(t, worker_address)
	}
}

// find waiting worker from workerstack channel
func findWorker() string {
	var worker_address string
	for {
		reply := make(chan string)
		msg := WorkerStackMsg{Type:"pop", Reply:reply}
		workerStackChan <- msg
		worker_address = <-reply

		if worker_address == "" {
			fmt.Println("no waiting workers")
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	return worker_address
}

func sendTask(task Task, worker_address string) {
	out, err := net.Dial("tcp", worker_address)
	encoder := gob.NewEncoder(out)
	err = encoder.Encode(task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("send task to %v : %v\n", worker_address, task)
}
