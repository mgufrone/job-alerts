// Code generated by mgufrone.dev/job-alerts/cmd/generate-domain, DO NOT EDIT.
package notification

import (
	uuid "github.com/google/uuid"
	try "github.com/mgufrone/go-utils/try"
	channel "mgufrone.dev/job-alerts/internal/domain/channel"
	job "mgufrone.dev/job-alerts/internal/domain/job"
	user "mgufrone.dev/job-alerts/internal/domain/user"
	"time"
)

func New(id uuid.UUID, user *user.Entity, job *job.Entity, channel *channel.Entity, isSent bool, sentAt time.Time, readAt *time.Time, createdAt time.Time, updatedAt time.Time) (ent *Entity, err error) {
	var res Entity
	err = try.RunOrError(func() error {
		return res.SetID(id)
	}, func() error {
		return res.SetUser(user)
	}, func() error {
		return res.SetJob(job)
	}, func() error {
		return res.SetChannel(channel)
	}, func() error {
		return res.SetIsSent(isSent)
	}, func() error {
		return res.SetSentAt(sentAt)
	}, func() error {
		return res.SetReadAt(readAt)
	}, func() error {
		return res.SetCreatedAt(createdAt)
	}, func() error {
		return res.SetUpdatedAt(updatedAt)
	})
	if err == nil {
		ent = &res
	}
	return
}
