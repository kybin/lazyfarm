package main

import (
	"fmt"
)

func removeDuplicates(a []int) []int {
        result := []int{}
        seen := map[int]int{}
        for _, val := range a {
                if _, ok := seen[val]; !ok {
                        result = append(result, val)
                        seen[val] = val
                }
        }
        return result
}

type intSlice []int

func (slice intSlice) pos(value int) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
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
