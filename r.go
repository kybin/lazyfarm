package main

type R struct {
	Run, Scene, Driver, Frame string
}

type Rseq struct {
	Run, Scene, Driver string
	Frames []int
}

