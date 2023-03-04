package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/db"
)

type Channel struct {
	db.Entity
	UserID      uuid.UUID `gorm:"index"`
	Name        string    `gorm:"index"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"index"`
	ScheduleAt  string
	Channels    pgtype.Array[string] `gorm:"type:text[]"`
	Criterias   string               `gorm:"type:text;"`
	db.EntityTimestamp
}

func (c *Channel) FromDomain(f *channel.Entity) {
	c.ID = f.ID()
	c.UserID = f.User().ID()
	c.Name = f.Name()
	c.Description = f.Description()
	c.IsActive = f.IsActive()
	c.ScheduleAt = f.ScheduleAt()
	var (
		chans pgtype.Array[string]
	)
	chans.Elements = f.Channels()
	c.Channels = chans
	c.Criterias = string(f.Criterias())
	c.CreatedAt = f.CreatedAt()
	c.UpdatedAt = f.UpdatedAt()
}

func (c *Channel) Transform() (*channel.Entity, error) {
	var (
		usr user.Entity
	)
	usr.SetID(c.UserID)
	return channel.New(
		c.ID,
		&usr,
		c.Name,
		c.Description,
		c.IsActive,
		c.ScheduleAt,
		c.CreatedAt,
		c.UpdatedAt,
		c.Channels.Elements,
		[]byte(c.Criterias),
	)
}
