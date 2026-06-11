package domain

import "time"

type Report struct {
	Success    bool
	ErrorTypes map[ErrorType]struct{}
	Time       time.Time
	Descriptor Descriptor
	Summary    Summary
}
