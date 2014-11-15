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
	go listen()

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

func handleWorker(worker *Worker) {
	var msgtype string

	switch worker.Status {
	case Login:
		msgtype = "push"
	case Logout:
		msgtype = "delete"
	default:
		log.Fatal("unknown status")
	}

	workerStackChan <- WorkerStackMsg{Type:msgtype, WorkerAddress:worker.Address}
}

func listen() {
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
		fmt.Println("data in")

		decoder := gob.NewDecoder(conn)
		var msgtype string
		err = decoder.Decode(&msgtype)
		if err != nil {
			log.Printf("warn : unknown data come in -")
			log.Println(err)
			continue
		}
		switch msgtype {
		case "job":
			fmt.Println("data is a Job")

			var j Job
			err = decoder.Decode(&j)
			if err != nil {
				log.Fatal(err)
			}

			go handleJob(&j)
		case "worker":
			fmt.Println("data is a Worker")

			var w Worker
			err = decoder.Decode(&w)
			if err != nil {
				log.Fatal(err)
			}

			go handleWorker(&w)
		default:
			log.Fatal(errors.New("Cannot determine data type"))
		}
	}
}

func handleJob(job *Job) {
	fmt.Printf("job added - %v\n", job)
	// separate it to tasks
	tasks, err := jobToTasks(job)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tasks {
		fmt.Printf("%v\n", t)
		worker_address := findWorker()
		go sendTask(t, worker_address, job.MaximumRetry)
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

func sendTask(task Task, worker_address string, maxretry int) {
	// todo : how to return log? with return statement or with channel?
	retry := 0
	for {
		// send task to worker
		out, err := net.Dial("tcp", worker_address)
		encoder := gob.NewEncoder(out)
		err = encoder.Encode(task)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("send task to %v : %v\n", worker_address, task)

		// wait for result
		var result Status

		decoder := gob.NewDecoder(out)
		err = decoder.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		// let our worker back
		workerStackChan <- WorkerStackMsg{Type:"push", WorkerAddress:worker_address}

		// see the result
		if result == Done {
			fmt.Println("task finished : ", worker_address)
			return
		} else {
			fmt.Println("task failed for some reason : ", worker_address)
			// check we can retry
			if retry >= maxretry {
				fmt.Println("exceeded max retry. task failed.")
				return
			} else {
				retry++
				fmt.Printf("retry %v\n", retry)
				worker_address = findWorker()
				continue
			}
		}
	}
}
