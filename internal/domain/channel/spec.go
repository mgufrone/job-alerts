package channel

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereOwner(id uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("user_id", criteria.Eq, id.String())
}
func WhereID(id uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("id", criteria.Eq, id.String())
}

func WhereActive() criteria.ICondition {
	return criteria.NewCondition("is_active", criteria.Eq, true)
}
func WhereInactive() criteria.ICondition {
	return criteria.NewCondition("is_active", criteria.Eq, false)
}

func HasChannel(channel string) criteria.ICondition {
	return criteria.NewCondition("channels", criteria.Like, channel)
}

func WhereCriteria(crit []byte) criteria.ICondition {
	return criteria.NewCondition("criterias", criteria.Eq, crit)
}
func WhereName(title string) criteria.ICondition {
	return criteria.NewCondition("title", criteria.Eq, title)
}
