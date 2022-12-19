package notification

import (
	"context"
	"errors"
	"github.com/mgufrone/go-utils/try"
	"mgufrone.dev/job-alerts/internal/domain/channel"
)

func (u *UseCase) Subscribe(ctx context.Context, ent *channel.Entity) (res error) {
	// check if exists
	var (
		q     = u.chQuery()
		cmd   = u.chCommand()
		count int64
	)
	res = try.RunOrError(func() (err error) {
		if ent == nil {
			err = errors.New("channel required")
		}
		return
	}, func() (err error) {
		cb := q.CriteriaBuilder()
		cb = cb.And(
			cb.Where(
				channel.WhereOwner(ent.User().ID()),
			),
			cb.Or(
				cb.Where(channel.WhereCriteria(ent.Criterias())),
				cb.Where(channel.WhereName(ent.Name())),
			),
		)
		count, err = q.Apply(cb).Count(ctx)
		return
	}, func() (err error) {
		if count > 0 {
			err = errors.New("channel exists")
		}
		return
	}, func() error {
		return cmd.Create(ctx, ent)
	})
	return
}
