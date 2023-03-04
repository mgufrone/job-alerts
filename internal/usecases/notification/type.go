package notification

import (
	"github.com/sirupsen/logrus"
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
	logger     logrus.FieldLogger
}

func NewUseCase(notifQuery notification.QueryResolver, notifCmd notification.CommandResolver, chQuery channel.QueryResolver, chCommand channel.CommandResolver, usrQuery user.QueryResolver, usrChQuery user_channel.QueryResolver, jQuery job.QueryResolver, publishers *publisher.Collection, logger logrus.FieldLogger) *UseCase {
	return &UseCase{notifQuery: notifQuery, notifCmd: notifCmd, chQuery: chQuery, chCommand: chCommand, usrQuery: usrQuery, usrChQuery: usrChQuery, jQuery: jQuery, publishers: publishers, logger: logger}
}
