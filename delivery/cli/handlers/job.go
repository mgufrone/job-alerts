package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/mgufrone/go-utils/try"
	"github.com/sirupsen/logrus"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/worker"
)

type PullParams struct {
	Worker worker.IWorker
	Debug  bool
	DryRun bool
}
type Job struct {
	logger logrus.FieldLogger
	query  job.QueryResolver
	cmd    job.CommandResolver
}

func NewJob(logger logrus.FieldLogger, query job.QueryResolver, cmd job.CommandResolver) *Job {
	return &Job{logger: logger, query: query, cmd: cmd}
}

func (j *Job) Pull(ctx context.Context, input PullParams) error {
	j.logger.Infoln("input", input.Debug, input.DryRun)
	var (
		jbs         []*job.Entity
		query       = j.query()
		cmd         = j.cmd()
		criterias   []criteria.ICriteriaBuilder
		counts      = map[string]int{}
		canContinue = true
	)
	err := try.RunOrError(func() (err error) {
		j.logger.Debug("fetching jobs")
		jbs, err = input.Worker.Fetch(ctx)
		j.logger.Debugf("found %d jobs", len(jbs))
		return
	}, func() error {
		for _, jb := range jbs {
			if jb.ID() == uuid.Nil {
				jb.SetID(uuid.New())
			}
			counts[jb.JobURL()] = 0
			criterias = append(criterias, query.CriteriaBuilder().
				Where(job.WhereURL(jb.JobURL()), job.WhereSource(jb.Source())))
			// filter out if it's already exists
		}
		return nil
	}, func() (err error) {
		crit := query.CriteriaBuilder().
			Or(criterias...)
		all, total, err := query.Apply(crit).GetAndCount(ctx)
		if err != nil {
			return
		}
		if int(total) == len(jbs) {
			canContinue = false
			return
		}
		for _, v := range all {
			counts[v.JobURL()] = 1
		}
		for k, v := range jbs {
			if counts[v.JobURL()] == 1 {
				continue
			}
			var err2 error
			jbs[k], err2 = input.Worker.FetchJob(ctx, v)
			if err2 == nil && jbs[k] == nil {
				jbs[k] = v
			}
			if err2 != nil {
				counts[v.JobURL()] = 1
				continue
			}
		}
		return
	})
	if err != nil {
		j.logger.Errorln("failed to proceed", err)
		return err
	}
	if !canContinue {
		j.logger.Debugln("nothing to insert. skipping")
		return nil
	}

	for _, jb := range jbs {
		if v, ok := counts[jb.JobURL()]; ok && v == 1 {
			continue
		}
		j.logger.Debugf("saving %s from %s to db", jb.JobURL(), jb.Source())
		if input.DryRun {
			continue
		}
		if err := cmd.Create(ctx, jb); err != nil {
			j.logger.Errorln("failed to save: ", jb.ID(), jb.JobURL(), err)
			continue
		}
	}
	return nil
}
