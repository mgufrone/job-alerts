package notification

import (
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/notification"
)

type UseCase struct {
	notifQuery notification.QueryResolver
	notifCmd   notification.CommandResolver

	chQuery   channel.QueryResolver
	chCommand channel.CommandResolver
}
