package main

import (
	"log"
	"net"
	"fmt"
	"encoding/gob"
	"time"
)

var workers = make(map[string]string)

func main() {
	go listenWorker()
	for {
		time.Sleep(time.Second)
		fmt.Println(workers)
	}
	// go listenJob()
}

func listenWorker() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connected")
		worker, status := handleWorkerConn(conn)
		if status == "waiting" {
			workers[worker.Address] = "waiting"
		} else if status == "logout" {
			delete(workers, worker.Address)
		} else {
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


func listenJob() {
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
		go handleJob(job)
	}
}

func handleJob(job *Job) {
	tasks := jobToTasks(job)
	for t := range tasks {
		worker_address := findingWorker()
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

func findingWorker() string {
	for {
		for w := range workers {
			if workers[w] == "waiting" {
				workers[w] = "processing"
				return w
			}
		}
		// O.K. We don't have any available worker. Wait a little.
		time.Sleep(time.Second)
	}
}
