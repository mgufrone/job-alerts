package channel

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type queryRepository struct {
	instance *db.Instance
	cb       criteria.ICriteriaBuilder
}

func NewQueryRepository(instance *db.Instance) channel.IQueryRepository {
	return &queryRepository{instance: instance}
}

func (q *queryRepository) CriteriaBuilder() criteria.ICriteriaBuilder {
	return db.NewCriteriaBuilder("channels")
}

func (q *queryRepository) Apply(cb criteria.ICriteriaBuilder) channel.IQueryRepository {
	q.cb = cb
	return q
}

func (q *queryRepository) GetAll(ctx context.Context) (out []*channel.Entity, err error) {
	defer q.resetCriteria()
	var jobs []*models.Channel
	source := q.run(ctx, q.cb)
	tx := source.Find(&jobs)
	err = tx.Error
	if err != nil {
		return
	}
	for _, ent := range jobs {
		tr, _ := ent.Transform()
		out = append(out, tr)
	}
	return
}

func (q *queryRepository) run(ctx context.Context, criteria criteria.ICriteriaBuilder) *gorm.DB {
	var mdl models.Channel
	if criteria == nil {
		return q.instance.Run(ctx, nil, &mdl)
	}
	d := q.instance.Run(ctx, criteria, &mdl)
	return d
}
func (q *queryRepository) Count(ctx context.Context) (total int64, err error) {
	defer q.resetCriteria()
	err = q.run(ctx, q.cb).Count(&total).Error
	return
}

func (q *queryRepository) GetAndCount(ctx context.Context) (out []*channel.Entity, total int64, err error) {
	cr := q.cb
	total, err = q.Apply(cr).Count(ctx)
	if err != nil || total == 0 {
		return
	}
	out, err = q.Apply(cr).GetAll(ctx)
	return
}

func (q *queryRepository) FindByID(ctx context.Context, id uuid.UUID) (out *channel.Entity, err error) {
	cr := q.CriteriaBuilder().Where(channel.WhereID(id))
	all, count, err := q.Apply(cr).GetAndCount(ctx)
	if err != nil {
		return
	}
	if count == 0 {
		err = fmt.Errorf("entity with ID: %s not found", id.String())
	}
	out = all[0]
	return
}

func (q *queryRepository) resetCriteria() {
	q.cb = nil
}
