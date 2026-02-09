package worker

import "time"

type Job struct {
	OrderID int64
	Timer   time.Duration
}

type JobInput struct {
	Job    Job
	Output chan JobOutput
}

type JobOutput struct {
	Job Job
	Err error
}
