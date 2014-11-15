package main

import (
	"time"
	_ "fmt"
)

type Task struct {
	Cmd string
	Frame int
}

type Status string

const (
	Wait Status = "wait"
	Processing Status = "processing"
	Done Status = "done"
	Failed Status = "failed"
)

type FrameInfo struct {
	Status Status
	Retry int
}

type Job struct {
	Cmd string
	Frames map[int]FrameInfo
	Submited time.Time
	Started time.Time
	Ended time.Time
	MaximumRetry int
	Broadcast bool
}

func (j *Job) NFrame() int {
	return len(j.Frames)
}

func (j *Job) NDone() int {
	n := 0
	for _, f := range j.Frames {
		if f.Status == Done {
			n++
		}
	}
	return n
}

// It checks how many frames in the job are failed. 
// It does not interested in how many retry attempts within a frame.
func (j *Job) NFail() int {
	n := 0
	for _, f := range j.Frames {
		if f.Retry != 0 {
			n++
		}
	}
	return n
}

type WorkerStatus string

const (
	Login WorkerStatus = "login"
	Logout WorkerStatus = "logout"
)

type Worker struct {
	Address string
	Status WorkerStatus
}

type WorkerStackMsg struct {
	Type string
	WorkerAddress string
	Reply chan string
}


// type Worker struct {
// 	Name string
// 	IP string
// 	Group string
// 	Exclusive bool
// }
