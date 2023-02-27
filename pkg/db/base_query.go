package db

import (
	"context"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

func (i *Instance) Run(ctx context.Context, crx criteria.ICriteriaBuilder, model interface{}) *gorm.DB {
	db2 := i.DB().WithContext(ctx)
	if crx != nil {
		cr := crx.(CriteriaBuilder)
		db2 = cr.Apply(db2)
	}
	return db2.Model(model)
}

func (i *Instance) Transaction(ctx context.Context, onRun func(tx *gorm.DB) error) error {
	return i.DB().WithContext(ctx).Transaction(onRun)
}
