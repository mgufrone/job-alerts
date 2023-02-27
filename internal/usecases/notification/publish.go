package notification

import (
	"context"
	"errors"
	"github.com/mgufrone/go-utils/try"
	"golang.org/x/sync/errgroup"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"sync"
)

func (u *UseCase) Publish(ctx context.Context, notif *notification.Entity) error {
	if notif == nil {
		return errors.New("cannot proceed")
	}
	var (
		chQuery    = u.chQuery()
		jQuery     = u.jQuery()
		usrQuery   = u.usrQuery()
		usrChQuery = u.usrChQuery()
		eg         errgroup.Group

		jobFound, userFound, channelFound bool
		channels                          []*user_channel.Entity
	)
	// check the existence of the data
	eg.Go(func() (err error) {
		cb := jQuery.CriteriaBuilder().Where(job.WhereID(notif.Job().ID()))
		count, err := jQuery.Apply(cb).Count(ctx)
		jobFound = count == 1
		return
	})
	eg.Go(func() (err error) {
		cb := usrQuery.CriteriaBuilder().
			Where(user.WhereID(notif.User().ID()))
		count, err := usrQuery.Apply(cb).Count(ctx)
		userFound = count == 1
		return
	})
	eg.Go(func() (err error) {
		cb := chQuery.CriteriaBuilder().
			Where(channel.WhereID(notif.Channel().ID()), channel.WhereActive())
		count, err := chQuery.Apply(cb).Count(ctx)
		channelFound = count == 1
		return
	})

	return try.RunOrError(func() error {
		return eg.Wait()
	}, func() error {
		if !(jobFound &&
			userFound &&
			channelFound) {
			// delete notif
			u.notifCmd().Delete(ctx, notif)
			return errors.New("notification cannot be sent due to requirements")
		}
		return nil
	}, func() (err error) {
		var count int64
		cb := usrChQuery.CriteriaBuilder()
		channelTypes := notif.Channel().Channels()
		cb.
			Where(user_channel.WhereOwner(notif.User())).
			Where(user_channel.WhereChannelTypes(channelTypes)).
			Where(user_channel.WhereActive())
		channels, count, err = usrChQuery.Apply(cb).GetAndCount(ctx)
		if count == 0 {
			return errors.New("channels not found. no need to proceed")
		}
		return
	}, func() error {
		var wg sync.WaitGroup

		for k, _ := range channels {
			wg.Add(1)
			go func(wg *sync.WaitGroup, channels []*user_channel.Entity, key int) {
				ch := channels[k]
				defer wg.Done()
				pub, err := u.publishers.Get(ch.ChannelType())
				if err == nil {
					pub.SetReceiver(ch)
					pub.Publish(ctx, notif)
				}
			}(&wg, channels, k)
		}
		wg.Wait()
		return nil
	})
}
