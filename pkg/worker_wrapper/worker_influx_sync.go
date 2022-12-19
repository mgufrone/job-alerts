package wrapper

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/common"
	"mgufrone.dev/job-alerts/pkg/worker"
)

const (
	perBatch = 100
)

type influxSync struct {
	source job.IQuery
	dest   job.ICommand
	logger *logrus.Entry
}

func NewInfluxSync(source job.IQuery, dest job.ICommand, logger *logrus.Entry) worker.IRunner {
	return &influxSync{source: source, dest: dest, logger: logger.WithField("component", "influx-sync")}
}

func (i *influxSync) getAndSend(ctx context.Context, page int) {
	cb := i.source.CriteriaBuilder()
	i.logger.Print("start batch at ", page, perBatch)
	jbs, err := i.source.Apply(cb.Paginate(page, perBatch)).GetAll(ctx)
	if err != nil {
		i.logger.WithField("batch", page).Error(err)
		return
	}
	for _, jb := range jbs {
		if err = i.dest.Create(ctx, jb); err != nil {
			i.logger.WithField("job", jb).Error(err)
		}
	}
}
func (i *influxSync) Run(ctx context.Context, args ...*job.Entity) error {
	var (
		total int64
	)
	err := common.TryOrError(func() (err error) {
		i.logger.Println("component", i.source, i.dest)
		total, err = i.source.Apply(i.source.CriteriaBuilder()).Count(ctx)
		return
	}, func() error {
		totalPages := math.Ceil(float64(total) / float64(perBatch))
		for idx := 0; idx < int(totalPages); idx++ {
			i.getAndSend(ctx, idx+1)
		}
		return nil
	})
	if err != nil {
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}
		if stackErr, ok := err.(stackTracer); ok {
			i.logger.WithField("stacktrace", stackErr.StackTrace()).Error()
		}
	}
	return err
}
