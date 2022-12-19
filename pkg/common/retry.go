package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultRetry = 3
	defaultDelay = time.Millisecond * 250
)

type ExecRetry func() error
type Retry struct {
	attempt  int
	maxRetry int
	delay    time.Duration
	exec     ExecRetry
	logger   *log.Entry
}

func DefaultRetry(retry ExecRetry, lg *log.Entry) *Retry {
	return NewRetry(defaultRetry, defaultDelay, retry, lg)
}

func NewRetry(maxRetry int, delay time.Duration, retry ExecRetry, lg *log.Entry) *Retry {
	return &Retry{
		attempt:  1,
		maxRetry: maxRetry,
		delay:    delay,
		exec:     retry,
		logger:   lg,
	}
}

func (r *Retry) Run() error {
	if err := r.exec(); err != nil {
		r.attempt++
		if r.attempt < r.maxRetry {
			r.logger.Warn(fmt.Sprintf("retry delayed at %s. retrying for %d out of %d", r.delay.String(), r.attempt, r.maxRetry))
			time.Sleep(r.delay * time.Duration(r.attempt))
			return r.Run()
		}
		return err
	}
	return nil
}
