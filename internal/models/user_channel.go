package models

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"mgufrone.dev/job-alerts/pkg/db"
)

type UserChannel struct {
	db.Entity
	UserID      uuid.UUID `gorm:"index"`
	ChannelType string    `gorm:"index"`
	Receiver    string    `gorm:"type:text"`
	IsActive    bool      `gorm:"index"`
	db.EntityTimestamp
}

func (u *UserChannel) FromDomain(ent *user_channel.Entity) {
	u.ID = ent.ID()
	u.UserID = ent.User().ID()
	u.ChannelType = ent.ChannelType()
	u.Receiver = ent.Receiver()
	u.IsActive = ent.IsActive()
	u.CreatedAt = ent.CreatedAt()
	u.UpdatedAt = ent.UpdatedAt()
}

func (u *UserChannel) Transform() (ent *user_channel.Entity, err error) {
	var (
		usr user.Entity
	)
	usr.SetID(u.UserID)
	return user_channel.New(
		u.ID,
		&usr,
		u.ChannelType,
		u.Receiver,
		u.IsActive,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
