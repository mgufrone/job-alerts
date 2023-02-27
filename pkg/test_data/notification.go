package test_data

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"time"
)

func ValidNotification(usr *user.Entity, job *job.Entity, ch *channel.Entity) *notification.Entity {
	nt, _ := notification.New(
		uuid.New(),
		usr,
		job,
		ch,
		false,
		time.Now(),
		nil,
		time.Now(),
		time.Now(),
	)
	return nt
}
