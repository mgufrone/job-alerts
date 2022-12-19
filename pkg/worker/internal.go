package worker

import (
	"context"
)

type Retryable struct {
}

func (r *Retryable) ReQueue(ctx context.Context, data []byte) error {
	panic("implement me")
}

func (r *Retryable) MaxRetry() int {
	panic("implement me")
}
