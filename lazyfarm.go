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

	groupmsgchan := make(chan GroupChanMsg)
	go groupChanMap(groupmsgchan)
	go createGroup("", groupmsgchan, msgchan, popchan) // default fallback group

	go listenWorker(msgchan)
	go listenJob(groupmsgchan)

	for {
		time.Sleep(time.Second)
		fmt.Println("tired...")
	}
}

func workerStack(msgchan chan WorkerStackMsg) {
	stack := make([]string, 0)
	for {
		msg := <-msgchan
		switch msg.Type {
		case "login": // push
			stack = append(stack, msg.WorkerAddress)
		case "logout": // find index and delete
			idx := -1
			for i, v := range stack {
				if (v == msg.WorkerAddress) {
					idx = i
					break
				}
			}
			if idx == -1 {
				notfound := errors.New("not found " + msg.WorkerAddress)
				log.Fatal(notfound)
			}
			stack = append(stack[:idx], stack[idx+1:]...) // delete address from stack
		case "waiting", "done": // same with login yet.
			stack = append(stack, msg.WorkerAddress)
		case "need": // pop
			if len(stack) == 0 {
				address = ""
			} else {
				last := len(stack)-1
				address = stack[last]
				stack = stack[:last]
			}
			msg.Reply <- address
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
		case "login", "logout", "done":
			fmt.Printf("%v, %v\n", worker, status)
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


func listenJob(groupchan chan GroupChanMsg) {
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
		fmt.Println("job decoded")
		go handleJob(job, groupchan)
	}
}

func handleJob(job *Job, groupfindchan chan GroupChanMsg) {
	fmt.Println("job will be handled")

	reply := make(chan chan Task)
	groupfindchan <- GroupChanMsg{Type:"find", GroupName:job.Group, Reply:reply}
	jobgroupchan := <-reply

	tasks := jobToTasks(job)
	fmt.Printf("%v\n", tasks)
	for _, t := range tasks {
		fmt.Printf("%v\n", t)
		jobgroupchan <- t
	}
}

func jobToTasks(job *Job) []Task {
	fmt.Println("job to tasks")
	nframes := len(job.Frames)
	tasks := make([]Task, nframes)
	for i := 0 ; i < nframes ; i++ {
		tasks[i].Run = job.Run
		tasks[i].Scene = job.Scene
		tasks[i].Driver = job.Driver
		tasks[i].Frame = job.Frames[i]
	}
	return tasks
}

func createGroup(name string, groupmsgchan chan GroupChanMsg, workermsgchan chan string, workerpopchan chan string) {
	grouptaskchan := make(chan Task)
	go handleGroupTask(grouptaskchan, workermsgchan, workerpopchan)

	groupmsgchan <- GroupChanMsg{Type:"add", GroupName:name, GroupChannel:grouptaskchan}
}

func handleGroupTask(grouptaskchan chan Task, workermsgchan chan string, workerpopchan chan string) {
	for {
		task := <-grouptaskchan
		worker_address := findWorker(workermsgchan, workerpopchan)
		sendTask(task, worker_address)
	}
}

func sendTask(task Task, worker_address string) {
	out, err := net.Dial("tcp", worker_address)
	encoder := gob.NewEncoder(out)
	err = encoder.Encode(task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("redirect : %v\n", task)
}

func groupChanMap(msgchan chan GroupChanMsg) {
	groupchanmap := make(map[string]chan Task)
	for {
		msg := <-msgchan
		switch msg.Type {
		case "add":
			groupchanmap[msg.GroupName] = msg.GroupChannel
		case "delete":
			delete(groupchanmap, msg.GroupName)
		case "find":
			msg.Reply <- groupchanmap[msg.GroupName]
		default:
			log.Fatal(errors.New(fmt.Sprintf("not expected message type '%v'", msg.Type)))
		}
		fmt.Println(groupchanmap)
	}
}

func findWorker(workermsgchan chan string, workerpopchan chan string) string {
	worker_address := ""
	for {
		workermsgchan <- "need"
		worker_address = <-workerpopchan
		if worker_address == "" {
			fmt.Println("no valid workers")
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}
	fmt.Printf("%v\n", worker_address)
	return worker_address
}
