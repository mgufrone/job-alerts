package notification

import (
	"context"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/pkg/db"
)

type commandRepository struct {
	db *db.Instance
}

func NewCommand(db *db.Instance) notification.ICommandRepository {
	return &commandRepository{db: db}
}

func (c *commandRepository) Create(ctx context.Context, in *notification.Entity) (err error) {
	var mdl models.Notification
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Create(&mdl).Error
	})
}

func (c *commandRepository) Update(ctx context.Context, in *notification.Entity) (err error) {
	var mdl models.Notification
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Updates(&mdl).Error
	})
}

func (c *commandRepository) Delete(ctx context.Context, in *notification.Entity) (err error) {
	var mdl models.Notification
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Delete(&mdl).Error
	})
}
