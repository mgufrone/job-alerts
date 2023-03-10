// Code generated by mgufrone.dev/job-alerts/cmd/generate-domain, DO NOT EDIT.
package channel

import (
	uuid "github.com/google/uuid"
	try "github.com/mgufrone/go-utils/try"
	user "mgufrone.dev/job-alerts/internal/domain/user"
	"time"
)

func New(id uuid.UUID, user *user.Entity, name string, description string, isActive bool, scheduleAt string, createdAt time.Time, updatedAt time.Time, channels []string, criterias []byte) (ent *Entity, err error) {
	var res Entity
	err = try.RunOrError(func() error {
		return res.SetID(id)
	}, func() error {
		return res.SetUser(user)
	}, func() error {
		return res.SetName(name)
	}, func() error {
		return res.SetDescription(description)
	}, func() error {
		return res.SetIsActive(isActive)
	}, func() error {
		return res.SetScheduleAt(scheduleAt)
	}, func() error {
		return res.SetCreatedAt(createdAt)
	}, func() error {
		return res.SetUpdatedAt(updatedAt)
	}, func() error {
		return res.SetChannels(channels)
	}, func() error {
		return res.SetCriterias(criterias)
	})
	if err == nil {
		ent = &res
	}
	return
}
