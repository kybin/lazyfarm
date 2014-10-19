package www

import (
	"time"
)

type JobInfo struct {
	Priority int // bigger number gets higher priority
	Name string
	Id string
	User string
}

type JobTime struct {
	Submit, Start, End time.Time
}

type JobStatus struct {
	Waiting, Active, Done, Err string
}


