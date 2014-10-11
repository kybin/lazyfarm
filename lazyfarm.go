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

func main() {

	groupinfochan := make(chan GroupInfoMsg)
	go groupInfoMap(groupinfochan)
	go createGroup("", groupinfochan) // default fallback group
	go createGroup("fx", groupinfochan) // default fallback group
	go createGroup("render", groupinfochan) // default fallback group

	go listenWorker(groupinfochan)
	go listenJob(groupinfochan)

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
			address := ""
			if len(stack) != 0 {
				last := len(stack)-1
				address = stack[last]
				stack = stack[:last]
			}
			msg.Reply <- address
		default:
			notexpect := errors.New(fmt.Sprintf("not expected message type '%v'", msg.Type))
			log.Fatal(notexpect)
		}
		fmt.Printf("%v\n", stack)
	}
}


func listenWorker(groupinfochan chan GroupInfoMsg) {
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

		reply := make(chan Group)
		groupinfochan <- GroupInfoMsg{Type:"find", GroupName:worker.Group, Reply:reply}
		group := <-reply

		switch status {
		case "login", "logout", "done":
			fmt.Printf("%v, %v\n", worker, status)
			workerinfo := WorkerStackMsg{Type:status, WorkerAddress:worker.Address}
			group.WorkerChannel <- workerinfo
		default:
			log.Fatal("unknown status")
		}
	}
}

func handleWorkerConn(conn net.Conn) (Worker, string)  {
	decoder := gob.NewDecoder(conn)
	w := Worker{}
	err := decoder.Decode(&w)
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


func listenJob(groupinfochan chan GroupInfoMsg) {
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
		go handleJob(job, groupinfochan)
	}
}

func handleJob(job *Job, groupinfochan chan GroupInfoMsg) {
	fmt.Println("job will be handled")

	reply := make(chan Group)
	groupinfochan <- GroupInfoMsg{Type:"find", GroupName:job.Group, Reply:reply}
	group := <-reply

	tasks := jobToTasks(job)
	fmt.Printf("%v\n", tasks)
	for _, t := range tasks {
		fmt.Printf("%v\n", t)
		group.TaskChannel <- t
	}
}

func createGroup(name string, groupinfochan chan GroupInfoMsg) {
	workerchan := make(chan WorkerStackMsg)
	go workerStack(workerchan)

	taskchan := make(chan Task)
	go handleGroupTask(taskchan, workerchan)

	group := Group{TaskChannel:taskchan, WorkerChannel:workerchan}

	groupinfochan <- GroupInfoMsg{Type:"add", GroupName:name, Group:group}
}

func handleGroupTask(taskchan chan Task, workerstackchan chan WorkerStackMsg) {
	for {
		task := <-taskchan
		worker_address := findWorker(workerstackchan)
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

func groupInfoMap(msgchan chan GroupInfoMsg) {
	groupchanmap := make(map[string]Group)
	for {
		msg := <-msgchan
		switch msg.Type {
		case "add":
			groupchanmap[msg.GroupName] = msg.Group
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

func findWorker(workerstackchan chan WorkerStackMsg) string {
	worker_address := ""
	for {
		reply := make(chan string)
		msg := WorkerStackMsg{Type:"need", Reply:reply}
		workerstackchan <- msg
		worker_address = <-reply
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
