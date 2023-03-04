package user_channel

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func WhereID(id uuid.UUID) criteria.ICondition {
	return criteria.NewCondition("user_id", criteria.Eq, id.String())
}
func WhereOwner(user *user.Entity) criteria.ICondition {
	return criteria.NewCondition("user_id", criteria.Eq, user.ID().String())
}
func WhereActive() criteria.ICondition {
	return criteria.NewCondition("is_active", criteria.Eq, true)
}
func WhereChannelType(channelType string) criteria.ICondition {
	return criteria.NewCondition("channel_type", criteria.Eq, channelType)
}
func WhereChannelTypes(channelType []string) criteria.ICondition {
	return criteria.NewCondition("channel_type", criteria.In, channelType)
}
