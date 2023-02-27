//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/user_channel mgufrone.dev/job-alerts
package user_channel

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"time"
)

type Entity struct {
	id          uuid.UUID
	user        *user.Entity
	channelType string
	receiver    string // telegram: chat id; email: email (duh); slack: webhook
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func (e *Entity) IsActive() bool {
	return e.isActive
}

func (e *Entity) SetIsActive(isActive bool) (err error) {
	e.isActive = isActive
	return
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

func (e *Entity) SetUser(userID *user.Entity) (err error) {
	e.user = userID
	return
}

func (e *Entity) ChannelType() string {
	return e.channelType
}

func (e *Entity) SetChannelType(channelType string) (err error) {
	e.channelType = channelType
	return
}

func (e *Entity) Receiver() string {
	return e.receiver
}

func (e *Entity) SetReceiver(receiver string) (err error) {
	e.receiver = receiver
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
