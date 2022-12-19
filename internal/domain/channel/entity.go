//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/channel mgufrone.dev/job-alerts
package channel

import (
	"fmt"
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/errors"
	"mgufrone.dev/job-alerts/pkg/helpers"
	"time"
)

type Entity struct {
	id          uuid.UUID
	user        *user.Entity
	name        string
	description string
	isActive    bool
	scheduleAt  string // cron notation
	createdAt   time.Time
	updatedAt   time.Time
	channels    []string
	criterias   []byte
}

func (e *Entity) Criterias() []byte {
	return e.criterias
}

func (e *Entity) SetCriterias(criterias []byte) (err error) {
	e.criterias = criterias
	return
}

func (e *Entity) User() *user.Entity {
	return e.user
}

func (e *Entity) SetUser(user *user.Entity) error {
	e.user = user
	return nil
}

func (e *Entity) ScheduleAt() string {
	return e.scheduleAt
}

func (e *Entity) SetScheduleAt(scheduleAt string) (err error) {
	if !helpers.IsValidCron(scheduleAt) {
		err = errors.FieldError("scheduleAt", fmt.Errorf("invalid cron notation: %s", scheduleAt))
		return
	}
	e.scheduleAt = scheduleAt
	return
}

func (e *Entity) ID() uuid.UUID {
	return e.id
}

func (e *Entity) SetID(id uuid.UUID) error {
	e.id = id
	return nil
}

func (e *Entity) Name() string {
	return e.name
}

func (e *Entity) SetName(name string) error {
	e.name = name
	return nil
}

func (e *Entity) Description() string {
	return e.description
}

func (e *Entity) SetDescription(description string) error {
	e.description = description
	return nil
}

func (e *Entity) IsActive() bool {
	return e.isActive
}

func (e *Entity) SetIsActive(isActive bool) error {
	e.isActive = isActive
	return nil
}

func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Entity) SetCreatedAt(createdAt time.Time) error {
	e.createdAt = createdAt
	return nil
}

func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Entity) SetUpdatedAt(updatedAt time.Time) error {
	e.updatedAt = updatedAt
	return nil
}

func (e *Entity) Channels() []string {
	return e.channels
}

func (e *Entity) SetChannels(channels []string) error {
	e.channels = channels
	return nil
}
