package wrapper

import (
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/common"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/worker"

	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

type WorkerIndex struct {
	worker worker.IWorker
	cmd    job.ICommand
	query  job.IQuery
	logger *log.Entry
}

type IndexOption func(index *WorkerIndex)

func WithLogger(entry *log.Entry) IndexOption {
	return func(index *WorkerIndex) {
		index.logger = entry
	}
}

func NewWorkerIndex(wrk worker.IWorker, cmd job.ICommand, query job.IQuery, options ...IndexOption) worker.IRunner {
	runner := &WorkerIndex{worker: wrk, cmd: cmd, query: query}
	for _, opt := range options {
		opt(runner)
	}
	return runner
}

func (idx *WorkerIndex) batchGet(ctx context.Context, jobs []*job.Entity) ([]*job.Entity, error) {
	cb := idx.query.CriteriaBuilder()
	ors := make([]criteria.ICriteriaBuilder, len(jobs))
	keyVal := make(map[string]*job.Entity, len(jobs))
	for id, u := range jobs {
		key := fmt.Sprintf("%s-%s", u.Source(), u.JobURL())
		keyVal[key] = u
		ors[id] = cb.And(cb.
			Where(job.WhereSource(u.Source())).
			Where(job.WhereURL(u.JobURL())))
	}
	cb = cb.Paginate(0, len(ors)).Or(ors...)
	var (
		res []*job.Entity
	)
	ret := common.NewRetry(5, time.Millisecond*250, func() (err error) {
		res, err = idx.query.Apply(cb).GetAll(ctx)
		return
	}, idx.log())
	if err := ret.Run(); err != nil {
		return nil, err
	}
	for _, u := range res {
		key := fmt.Sprintf("%s-%s", u.Source(), u.JobURL())
		delete(keyVal, key)
	}
	var out []*job.Entity
	for _, j := range keyVal {
		if j == nil {
			continue
		}
		out = append(out, j)
	}
	return out, nil
}

func (idx *WorkerIndex) batch(ctx context.Context, jobs []*job.Entity, batchNumber int) []*job.Entity {
	batchLength := int(math.Ceil(float64(len(jobs)) / float64(batchNumber)))
	var (
		res []*job.Entity
	)
	for id := 0; id < batchLength; id++ {
		min := id * batchNumber
		max := (id + 1) * batchNumber
		if max > len(jobs) {
			max = len(jobs)
		}
		idx.log().Println("running batch at", id, min, max)
		jbs, err := idx.batchGet(ctx, jobs[min:max])
		if err != nil {
			idx.log().Println("error occurred", min, max)
			continue
		}
		if len(jbs) > 0 {
			res = append(res, jbs...)
		}
	}
	return res
}

func (idx *WorkerIndex) log() *log.Entry {
	if idx.logger != nil {
		return idx.logger
	}
	return log.NewEntry(log.StandardLogger())
}

func (idx *WorkerIndex) Run(ctx context.Context, _ ...*job.Entity) error {
	urls, err := idx.worker.Fetch(ctx)
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		return nil
	}
	idx.log().Println("processing jobs: ", len(urls))
	jbs := idx.batch(ctx, urls, 50)
	idx.log().Printf("queuing %d jobs", len(jbs))
	for _, j := range jbs {
		if err := idx.cmd.Create(ctx, j); err != nil {
			idx.log().WithFields(log.Fields{
				"source": j.Source(),
				"url":    j.JobURL(),
			}).Error(err)
		}
	}
	return nil
}
