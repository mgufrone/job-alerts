package worker

import (
	"context"
	"mgufrone.dev/job-alerts/internal/domain/job"
	http2 "net/http"
)

type IWorker interface {
	Fetch(ctx context.Context) ([]*job.Entity, error)
	FetchJob(ctx context.Context, job2 *job.Entity) (*job.Entity, error)
}

type IHTTPClient interface {
	Do(req *http2.Request) (*http2.Response, error)
	ToHTTPClient() *http2.Client
}

type IRunner interface {
	Run(ctx context.Context, args ...*job.Entity) error
}

type IRetryable interface {
	ReQueue(ctx context.Context, data []byte) error
	MaxRetry() int
	Failed(ctx context.Context, data []byte) error
}
