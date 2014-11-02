package main

import (
	"time"
)

type Task struct {
	Cmd string
	Frame int
}

type Status int

const (
	Wait Status = iota
	Processing
	Done
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

type Worker struct {
	Address string
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
