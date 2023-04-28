package app

import (
	"mgufrone.dev/job-alerts/delivery/cli/handlers"
)

type Kernel struct {
	Job *handlers.Job
}

func NewKernel(job *handlers.Job) *Kernel {
	return &Kernel{Job: job}
}
