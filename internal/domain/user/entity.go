//go:generate go run mgufrone.dev/job-alerts/cmd/generate-domain mgufrone.dev/job-alerts/internal/domain/user mgufrone.dev/job-alerts
package user

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	errors2 "mgufrone.dev/job-alerts/pkg/errors"
	"time"
)

type Status int
type Role int

const (
	Guest Role = 1 + iota<<16
	User
	Admin
	Premium
)

const (
	Deactivated Status = -1
	Active             = 1
)

type Entity struct {
	id        uuid.UUID
	authID    string
	status    Status
	roles     Role
	createdAt time.Time
	updatedAt time.Time
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

func (e *Entity) Roles() Role {
	return e.roles
}

func (e *Entity) SetRoles(roles Role) (err error) {
	if !(hasRoles(roles, Admin) || hasRoles(roles, User) || hasRoles(roles, Guest) || hasRoles(roles, Premium)) {
		return errors2.FieldError("roles", errors.New("invalid role"))
	}
	e.roles = roles
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
func hasRoles(roles Role, check Role) bool {
	return check&roles == check
}
func (e *Entity) HasRole(role Role) bool {
	if role == 0 {
		return false
	}
	return hasRoles(e.roles, role)
}

func (e *Entity) Copy() (*Entity, error) {
	return New(
		e.ID(),
		e.AuthID(),
		e.Status(),
		e.Roles(),
		e.CreatedAt(),
		e.UpdatedAt(),
	)
}
