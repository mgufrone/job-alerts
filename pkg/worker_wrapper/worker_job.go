package wrapper

import (
	"context"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/common"
	err2 "mgufrone.dev/job-alerts/pkg/errors"
	"mgufrone.dev/job-alerts/pkg/worker"
	"strings"
)

type WorkerResolver func(workerName string) worker.IWorker
type WorkerJob struct {
	worker WorkerResolver
	cmd    job.ICommandRepository
	query  job.IQueryRepository
	lg     *log.Entry
}

func NewWorkerJob(wrk WorkerResolver, query job.IQueryRepository, cmd job.ICommandRepository, lg *log.Entry) worker.IRunner {
	return &WorkerJob{wrk, cmd, query, lg}
}

func (wrk *WorkerJob) log() *log.Entry {
	if wrk.lg != nil {
		return wrk.lg
	}
	return log.NewEntry(log.StandardLogger())
}

func (wrk *WorkerJob) Run(ctx context.Context, arg ...*job.Entity) (err error) {
	var (
		jb          *job.Entity
		shouldBreak bool
	)
	return common.TryOrError(func() (err error) {
		jb = arg[0]
		cb := wrk.query.CriteriaBuilder().
			Where(
				job.WhereSource(jb.Source()),
				job.WhereURL(jb.JobURL()),
			)
		var (
			total int64
		)
		ret := common.DefaultRetry(func() (err error) {
			total, err = wrk.query.Apply(cb).Count(ctx)
			return
		}, wrk.log())
		if err = ret.Run(); err != nil {
			return
		}
		shouldBreak = total > 0
		return
	}, func() (err error) {
		if shouldBreak {
			return
		}
		msg := wrk.log().WithField("job", jb)
		msg.Printf("start executing job %s:%s\n", jb.Source(), jb.JobURL())
		runner := wrk.worker(jb.Source())
		if jb.ID() == uuid.Nil {
			_ = jb.SetID(uuid.New())
		}
		jb, err = runner.FetchJob(ctx, jb)
		if err != nil {
			msg.Error(err)
			var httpErr *err2.ClientError
			if errors.As(err, &httpErr) && strings.Contains(httpErr.Error(), "status: 429") {
				return nil
			}
			return
		}
		msg.Println("done")
		return wrk.cmd.Create(ctx, jb)
	})
}
