package main

type Task struct {
	Run, Scene, Driver, Frame string
}

type Job struct {
	Run, Scene, Driver string
	Frames []int
}

type Worker struct {
	Address string
}

// type Worker struct {
// 	Name string
// 	IP string
// 	Group string
// 	Exclusive bool
// }
