package notification

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/mgufrone/go-utils/try"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/helpers"
	job2 "mgufrone.dev/job-alerts/internal/repositories/job"
	"sync"
	"time"
)

func (u *UseCase) broadcast(ctx context.Context, ch *channel.Entity) error {
	// check if the owner is active
	var (
		usrQuery   = u.usrQuery()
		jQuery     = u.jQuery()
		notifQuery = u.notifQuery()
		jobs       []*job.Entity
		jobsCount  int64
	)
	return try.RunOrError(func() error {
		cb := usrQuery.CriteriaBuilder().Where(
			user.WhereActive(),
			user.WhereID(ch.User().ID()),
		)
		usrCount, err := usrQuery.Apply(cb).Count(ctx)
		if err != nil {
			return err
		}
		if usrCount == 0 {
			return errors.New("user is inactive, skipping")
		}
		return nil
	}, func() (err error) {
		// query jobs
		cb := helpers.FilterToCriteria(ch, jQuery.CriteriaBuilder(), job2.TagCriteria())
		// filter only job created one month prior
		cb = cb.Where(job.WhereBefore(time.Hour * 24 * 30))
		jobs, jobsCount, err = jQuery.Apply(cb).GetAndCount(ctx)
		if err == nil && jobsCount == 0 {
			err = errors.New("no jobs found, skipping")
		}
		return
	}, func() error {
		var (
			mappedJobs = map[uuid.UUID]*job.Entity{}
			compacted  []*job.Entity
		)
		for _, v := range jobs {
			mappedJobs[v.ID()] = v
		}
		//for k, v := range jobs {
		cb := notifQuery.CriteriaBuilder().
			Where(notification.WhereJobs(jobs...)).
			Where(notification.WhereOwner(ch.User()))

		all, count, err := notifQuery.Apply(cb).GetAndCount(ctx)
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		for _, jb := range all {
			if _, ok := mappedJobs[jb.Job().ID()]; ok {
				delete(mappedJobs, jb.Job().ID())
			}
			if len(mappedJobs) == 0 {
				break
			}
		}
		//}
		if len(mappedJobs) == 0 {
			return errors.New("found jobs already scheduled, skipping")
		}
		for _, jb := range mappedJobs {
			compacted = append(compacted, jb)
		}
		jobs = compacted
		return nil
	}, func() (err error) {
		schedule := time.Now()
		if ch.ScheduleAt() != "now" {
			crn := cronexpr.MustParse(ch.ScheduleAt())
			schedule = crn.Next(time.Now())
		}
		for _, j := range jobs {
			notif, err := notification.New(
				uuid.New(),
				ch.User(),
				j,
				ch,
				false,
				schedule,
				nil,
				time.Now(),
				time.Now(),
			)
			if err != nil {
				continue
			}
			if err = u.notifCmd().Create(ctx, notif); err != nil {
				continue
			}
		}
		return nil
	})
}
func (u *UseCase) Loop(ctx context.Context) error {
	var (
		chQuery       = u.chQuery()
		channels      []*channel.Entity
		channelsCount int64
	)
	return try.RunOrError(func() (err error) {
		cr := chQuery.CriteriaBuilder()
		cr.Where(channel.WhereActive())
		channels, channelsCount, err = chQuery.Apply(cr).GetAndCount(ctx)
		if err == nil && channelsCount == 0 {
			return errors.New("no channels found. skipping")
		}
		return
	}, func() (err error) {
		var wg sync.WaitGroup
		wg.Add(len(channels))
		for _, ch := range channels {
			go func(w *sync.WaitGroup, ch *channel.Entity) {
				defer w.Done()
				if err := u.broadcast(context.Background(), ch); err != nil {
					// do something about the error
					u.logger.Errorln(err)
				}
			}(&wg, ch)
		}
		wg.Wait()
		return nil
	})
}
