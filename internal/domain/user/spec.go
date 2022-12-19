package user

import "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"

func WhereAuthID(authID string) criteria.ICondition {
	return criteria.NewCondition("auth_id", criteria.Eq, authID)
}
