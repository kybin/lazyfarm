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
	go createGroup("fx", groupinfochan) // test group
	go createGroup("render", groupinfochan) // test group

	go listenWorker(groupinfochan)
	go listenJob(groupinfochan)

	for {
		time.Sleep(time.Second)
		// fmt.Println("tired...")
	}
}

func workerStack(groupname string, msgchan chan WorkerStackMsg) {
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
		fmt.Printf("group %v's workers - %v\n", groupname, stack)
	}
}


func listenWorker(groupinfochan chan GroupInfoMsg) {
	// worker socket
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

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

		handleWorker(status, worker, groupinfochan)
	}
}

func handleWorker(status string, worker Worker, groupinfochan chan GroupInfoMsg) {
	// find worker group
	reply := make(chan Group)
	groupinfochan <- GroupInfoMsg{Type:"find", GroupName:worker.Group, Reply:reply}
	group := <-reply

	// send worker status to the group's workerstack
	switch status {
	case "login", "logout", "done":
		workerinfo := WorkerStackMsg{Type:status, WorkerAddress:worker.Address}
		group.WorkerChannel <- workerinfo
	default:
		log.Fatal("unknown status")
	}
}

func listenJob(groupinfochan chan GroupInfoMsg) {
	// job socket
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

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

		go handleJob(job, groupinfochan)
	}
}

func handleJob(job *Job, groupinfochan chan GroupInfoMsg) {
	// find job's group
	reply := make(chan Group)
	groupinfochan <- GroupInfoMsg{Type:"find", GroupName:job.Group, Reply:reply}
	group := <-reply

	// separate it to tasks
	tasks := jobToTasks(job)

	for _, t := range tasks {
		fmt.Printf("%v\n", t)
		// the task will served by worker
		group.TaskChannel <- t
	}
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
	}
}

func createGroup(name string, groupinfochan chan GroupInfoMsg) {
	// create group's worker channel and task channel
	workerchan := make(chan WorkerStackMsg)
	go workerStack(name, workerchan)
	taskchan := make(chan Task)
	go handleGroupTask(taskchan, workerchan)
	// save it to group info map
	group := Group{TaskChannel:taskchan, WorkerChannel:workerchan}
	groupinfochan <- GroupInfoMsg{Type:"add", GroupName:name, Group:group}
	fmt.Printf("group '%v' created\n", name)
}

func handleGroupTask(taskchan chan Task, workerstackchan chan WorkerStackMsg) {
	for {
		task := <-taskchan
		worker_address := findWorker(workerstackchan)
		sendTask(task, worker_address)
	}
}

// find waiting worker from group's workerstack channel
func findWorker(workerstackchan chan WorkerStackMsg) string {
	var worker_address string
	for {
		reply := make(chan string)
		msg := WorkerStackMsg{Type:"need", Reply:reply}
		workerstackchan <- msg
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
