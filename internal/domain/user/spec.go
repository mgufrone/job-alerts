package user

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereAuthID(authID string) criteria.ICondition {
	return criteria.NewCondition("auth_id", criteria.Eq, authID)
}
func WhereID(id uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("id", criteria.Eq, id.String())
}
func WhereIDs(ids []uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("id", criteria.In, ids)
}
func WhereActive() criteria.ICondition {
	return criteria.NewCondition("status", criteria.Eq, Active)
}
