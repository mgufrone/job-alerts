package notification

import (
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereAuthID(authID string) criteria.ICondition {
	return criteria.NewCondition("auth_id", criteria.Eq, authID)
}

func WhereJob(job *job.Entity) criteria.ICondition {
	return criteria.NewCondition("job_id", criteria.Eq, job.ID().String())
}
func WhereOwner(job *user.Entity) criteria.ICondition {
	return criteria.NewCondition("user_id", criteria.Eq, job.ID().String())
}
