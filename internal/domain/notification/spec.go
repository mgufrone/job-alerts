package notification

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereJob(job *job.Entity) criteria.ICondition {
	return criteria.NewCondition("job_id", criteria.Eq, job.ID().String())
}
func WhereID(id uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("id", criteria.Eq, id.String())
}
func WhereJobs(job ...*job.Entity) criteria.ICondition {
	if len(job) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, len(job))
	for k, v := range job {
		ids[k] = v.ID()
	}
	return criteria.NewCondition("job_id", criteria.In, ids)
}
func WhereOwner(job *user.Entity) criteria.ICondition {
	return criteria.NewCondition("user_id", criteria.Eq, job.ID().String())
}
