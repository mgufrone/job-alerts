package wrapper

import (
	"context"
	"errors"
	"github.com/mgufrone/go-utils/try"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/worker"
)

type WorkerDataJob struct {
	main  job.ICommandRepository
	query job.IQueryRepository
}

func NewWorkerDataJob(main job.ICommandRepository, query job.IQueryRepository) worker.IRunner {
	return &WorkerDataJob{main: main, query: query}
}

func (w *WorkerDataJob) Run(ctx context.Context, args ...*job.Entity) (err error) {
	var count int64
	v := args[0]
	if v == nil {
		return errors.New("job should not be nil")
	}
	return try.RunOrError(func() (err error) {
		cb := w.query.CriteriaBuilder().
			Where(job.WhereSource(v.Source())).
			Where(job.WhereURL(v.JobURL()))
		count, err = w.query.Apply(cb).Count(ctx)
		return
	}, func() (err error) {
		if count != 0 {
			return nil
		}
		err = w.main.Create(ctx, v)
		if err != nil {
			if codeErr, ok := status.FromError(err); ok && codeErr.Code() == codes.AlreadyExists {
				return nil
			}
			log.Error("error saving data", err, v)
		}
		return
	})
}
