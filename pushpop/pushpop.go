package main

import (
	"time"
	"math/rand"
	"strings"
	"fmt"
	"strconv"
)

func main() {
	msgchan := make(chan string)
	popchan := make(chan string)
	go pushPopStack(msgchan, popchan)
	go pusher(msgchan)
	go poper(msgchan, popchan)
	time.Sleep(10*time.Second)
}

func pushPopStack(msgchan chan string, popchan chan string) {
	stack := make([]string, 0)
	var msg string
	var x string
	for {
		msg = <-msgchan
		if strings.HasPrefix(msg, "push") {
			msglist := strings.Split(msg, " ") // first is "push" string it self
			x = msglist[1]
			stack = append(stack, x)
			fmt.Printf("pushed : ")
		} else if msg == "pop" {
			if len(stack) == 0 {
				x = ""
			} else {
				last := len(stack)-1
				x = stack[last]
				stack = stack[:last]
			}
			popchan <- x
			fmt.Printf("Poped : ")
		} else {
			fmt.Printf("not expected (%v) ", msg)
		}
		fmt.Printf("%v\n", stack)
	}
}

func pusher(msgchan chan string) {
	var msg string
	var x string
	cmd := "push "
	for {
		time.Sleep(2*time.Duration(rand.Intn(1e2))*time.Millisecond)
		x = strconv.Itoa(rand.Intn(1000))
		msg = cmd+x
		msgchan <- msg
	}
}

func poper(msgchan chan string, popchan chan string) {
	for {
		time.Sleep(2*time.Duration(rand.Intn(1e2))*time.Millisecond)
		msgchan <- "pop"
		<-popchan // do not use it yet
	}
}



