package job

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type queryRepository struct {
	instance *db.Instance
	cb       criteria.ICriteriaBuilder
}

func NewQueryRepository(instance *db.Instance) job.IQueryRepository {
	return &queryRepository{instance: instance}
}

func (q *queryRepository) CriteriaBuilder() criteria.ICriteriaBuilder {
	return db.NewCriteriaBuilder("jobs")
}

func (q *queryRepository) Apply(cb criteria.ICriteriaBuilder) job.IQueryRepository {
	q.cb = cb
	return q
}

func (q *queryRepository) GetAll(ctx context.Context) (out []*job.Entity, err error) {
	defer q.resetCriteria()
	var jobs []*models.Job
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
	var mdl models.Job
	if criteria == nil {
		return q.instance.Run(ctx, nil, &mdl)
	}
	dbCr := criteria.(db.CriteriaBuilder)
	d := q.instance.Run(ctx, criteria, &mdl)
	if dbCr.Has("Tags") || dbCr.Prefix() == "Tags" {
		d = d.Joins("LEFT JOIN job_tags ON job_tags.job_id = jobs.id").
			Joins("LEFT JOIN tags Tags ON Tags.id = job_tags.tag_id")
	}
	return d
}
func (q *queryRepository) Count(ctx context.Context) (total int64, err error) {
	defer q.resetCriteria()
	err = q.run(ctx, q.cb).Count(&total).Error
	return
}

func (q *queryRepository) GetAndCount(ctx context.Context) (out []*job.Entity, total int64, err error) {
	cr := q.cb
	total, err = q.Apply(cr).Count(ctx)
	if err != nil || total == 0 {
		return
	}
	out, err = q.Apply(cr).GetAll(ctx)
	return
}

func (q *queryRepository) FindByID(ctx context.Context, id uuid.UUID) (out *job.Entity, err error) {
	cr := q.CriteriaBuilder().Where(job.WhereID(id))
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
