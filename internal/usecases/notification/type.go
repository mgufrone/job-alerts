package notification

import (
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"mgufrone.dev/job-alerts/internal/services/publisher"
)

type UseCase struct {
	notifQuery notification.QueryResolver
	notifCmd   notification.CommandResolver
	chQuery    channel.QueryResolver
	chCommand  channel.CommandResolver
	usrQuery   user.QueryResolver
	usrChQuery user_channel.QueryResolver
	jQuery     job.QueryResolver

	publishers *publisher.Collection
}
