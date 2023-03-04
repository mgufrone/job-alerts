package models

import (
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/db"
)

type User struct {
	db.Entity
	Roles  user.Role
	AuthID string      `gorm:"index"`
	Status user.Status `gorm:"index"`
	db.EntityTimestamp
}

func (u *User) FromDomain(ent *user.Entity) {
	u.ID = ent.ID()
	u.Status = ent.Status()
	u.Roles = ent.Roles()
	u.AuthID = ent.AuthID()
	u.CreatedAt = ent.CreatedAt()
	u.UpdatedAt = ent.UpdatedAt()
}
func (u *User) Transform() (*user.Entity, error) {
	return user.New(
		u.ID,
		u.AuthID,
		u.Status,
		u.Roles,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
