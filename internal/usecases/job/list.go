package job

import (
	"context"
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	job2 "mgufrone.dev/job-alerts/internal/repositories/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func (u *UseCase) listJobsForUser(ctx context.Context, usr *user.Entity) (jobs []*job.Entity, err error) {
	// only fetch last 100
	var (
		q = u.ntQuery()
	)
	cb := q.CriteriaBuilder().Where(notification.WhereOwner(usr)).Paginate(0, 100)
	notifs, err := q.Apply(cb).GetAll(ctx)
	if err != nil {
		return
	}
	for _, notif := range notifs {
		jobs = append(jobs, notif.Job())
	}
	return
}
func (u *UseCase) List(ctx context.Context, input ListInput) (jobs []*job.Entity, total int64, err error) {
	var (
		q      = u.jbQuery()
		wheres []criteria.ICriteriaBuilder
		jobIDs []*job.Entity
		cb     = q.CriteriaBuilder().Order("updated_at", "desc").Select(input.Fields...)
	)
	if input.User != nil {
		jobIDs, err = u.listJobsForUser(ctx, input.User)
		if err != nil {
			return nil, 0, err
		}
		if len(jobIDs) == 0 {
			return nil, 0, nil
		}
		ids := make([]uuid.UUID, len(jobIDs))
		for k, jb := range jobIDs {
			ids[k] = jb.ID()
		}
		wheres = append(wheres, cb.Where(job.WhereInID(ids)))
	}

	// force up to 50 if it's 0 or beyond 50
	if input.PerPage <= 0 || input.PerPage >= 50 {
		input.PerPage = 50
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Keyword != "" {
		wheres = append(wheres, cb.Or(
			cb.Where(job.WhereDescriptionContains(input.Keyword)),
			cb.Where(job.WhereRoleContains(input.Keyword)),
		))
	}
	if len(input.Skills) > 0 {
		wheres = append(wheres, job2.TagCriteria().Where(job.WhereTags(input.Skills)))
	}

	return q.Apply(cb.And(wheres...).Paginate(input.Page, input.PerPage)).GetAndCount(ctx)
}
