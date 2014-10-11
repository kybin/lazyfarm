package main

type Task struct {
	Run, Scene, Driver string
	Frame int
}

type Job struct {
	Run, Scene, Driver string
	Frames []int
	Group string
}

type Worker struct {
	Address string
	Group string
}

type WorkerStackMsg struct {
	Type string
	WorkerAddress string
	Reply chan string
}

type Group struct {
	TaskChannel chan Task
	WorkerChannel chan WorkerStackMsg
}

type GroupInfoMsg struct {
	Type string
	GroupName string
	Group Group
	Reply chan Group
}

// type Worker struct {
// 	Name string
// 	IP string
// 	Group string
// 	Exclusive bool
// }
