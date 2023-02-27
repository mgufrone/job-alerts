package job

import (
	"context"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/pkg/db"
)

type commandRepository struct {
	db *db.Instance
}

func NewCommand(db *db.Instance) job.ICommandRepository {
	return &commandRepository{db: db}
}

func (c *commandRepository) Create(ctx context.Context, in *job.Entity) (err error) {
	var mdl models.Job
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Create(&mdl).Error
	})
}

func (c *commandRepository) Update(ctx context.Context, in *job.Entity) (err error) {
	var mdl models.Job
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Updates(&mdl).Error
	})
}

func (c *commandRepository) Delete(ctx context.Context, in *job.Entity) (err error) {
	var mdl models.Job
	mdl.FromDomain(in)
	return c.db.Transaction(ctx, func(tx *gorm.DB) error {
		return tx.Delete(&mdl).Error
	})
}
