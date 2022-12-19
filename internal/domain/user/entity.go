//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/user mgufrone.dev/job-alerts
package user

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/pkg/helpers"
)

type Status int

const (
	Deactivated Status = -1
	Active             = 1
)

type Entity struct {
	id     uuid.UUID
	authID string
	status Status
	roles  []string
}

func (e *Entity) Roles() []string {
	return e.roles
}

func (e *Entity) SetRoles(roles []string) (err error) {
	if err = validation.Validate(
		roles,
		validation.Required,
		validation.Length(1, 0),
		validation.Each(
			validation.Required,
		),
	); err == nil {
		e.roles = roles
	}
	return
}

func (e *Entity) ID() uuid.UUID {
	return e.id
}

func (e *Entity) SetID(id uuid.UUID) (err error) {
	e.id = id
	return
}

func (e *Entity) AuthID() string {
	return e.authID
}

func (e *Entity) SetAuthID(authID string) error {
	if err := validation.Validate(authID, validation.Required); err != nil {
		return err
	}
	e.authID = authID
	return nil
}

func (e *Entity) Status() Status {
	return e.status
}

func (e *Entity) SetStatus(status Status) (err error) {
	if !(status == Active || status == Deactivated) {
		err = errors.New("invalid value")
		return
	}
	e.status = status
	return
}
func (e *Entity) HasRole(role string) bool {
	return helpers.Contains(e.roles, role)
}

func (e *Entity) Copy() (*Entity, error) {
	return New(
		e.ID(),
		e.AuthID(),
		e.Status(),
		e.Roles(),
	)
}
