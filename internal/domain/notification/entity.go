//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/notification mgufrone.dev/job-alerts
package notification

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"time"
)

type Entity struct {
	id        uuid.UUID
	user      *user.Entity
	job       *job.Entity
	isSent    bool
	sentAt    time.Time
	readAt    *time.Time
	createdAt time.Time
	updatedAt time.Time
}

func (e *Entity) ID() uuid.UUID {
	return e.id
}

func (e *Entity) SetID(id uuid.UUID) (err error) {
	e.id = id
	return
}

func (e *Entity) User() *user.Entity {
	return e.user
}

func (e *Entity) SetUser(user *user.Entity) (err error) {
	e.user = user
	return
}

func (e *Entity) Job() *job.Entity {
	return e.job
}

func (e *Entity) SetJob(job *job.Entity) (err error) {
	e.job = job
	return
}

func (e *Entity) IsSent() bool {
	return e.isSent
}

func (e *Entity) SetIsSent(isSent bool) (err error) {
	e.isSent = isSent
	return
}

func (e *Entity) SentAt() time.Time {
	return e.sentAt
}

func (e *Entity) SetSentAt(sentAt time.Time) (err error) {
	e.sentAt = sentAt
	return
}

func (e *Entity) ReadAt() *time.Time {
	return e.readAt
}

func (e *Entity) SetReadAt(readAt *time.Time) (err error) {
	e.readAt = readAt
	return
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) SetCreatedAt(createdAt time.Time) (err error) {
	e.createdAt = createdAt
	return
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) SetUpdatedAt(updatedAt time.Time) (err error) {
	e.updatedAt = updatedAt
	return
}
