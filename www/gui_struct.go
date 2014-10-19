package www

import (
	"time"
)

type WWWJobInfo struct {
	priority int // bigger number gets higher priority
	name string
	id string
	user string
}

type WWWJobTime struct {
	submit, start, end time.Time
}

type WWWJobStatus struct {
	waiting, active, done, err string
}


