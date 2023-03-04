package models

import (
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	user2 "mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/db"
	"time"
)

type Notification struct {
	db.Entity
	UserID    uuid.UUID `gorm:"index"`
	ChannelID uuid.UUID `gorm:"index"`
	JobID     uuid.UUID `gorm:"index"`
	IsSent    bool      `gorm:"index"`
	SentAt    time.Time
	ReadAt    *time.Time
	db.EntityTimestamp
}

func (n *Notification) FromDomain(ent *notification.Entity) {
	n.ID = ent.ID()
	n.UserID = ent.User().ID()
	n.ChannelID = ent.Channel().ID()
	n.JobID = ent.Job().ID()
	n.IsSent = ent.IsSent()
	n.SentAt = ent.SentAt()
	n.ReadAt = ent.ReadAt()
	n.UpdatedAt = ent.UpdatedAt()
	n.CreatedAt = ent.CreatedAt()
}

func (n *Notification) Transform() (*notification.Entity, error) {
	var (
		ch   channel.Entity
		user user2.Entity
		jb   job.Entity
	)
	user.SetID(n.UserID)
	ch.SetID(n.ChannelID)
	jb.SetID(n.JobID)
	return notification.New(
		n.ID,
		&user,
		&jb,
		&ch,
		n.IsSent,
		n.SentAt,
		n.ReadAt,
		n.CreatedAt,
		n.UpdatedAt,
	)
}
