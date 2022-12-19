package job

import (
	"context"
	criteria2 "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type IQuery interface {
	CriteriaBuilder() criteria2.ICriteriaBuilder
	Apply(cb criteria2.ICriteriaBuilder) IQuery
	// expose-able
	GetAll(ctx context.Context) (out []*Entity, err error)
	Count(ctx context.Context) (total int64, err error)
	GetAndCount(ctx context.Context) (out []*Entity, total int64, err error)
	FindByID(ctx context.Context, id string) (out *Entity, err error)
}

type ICommand interface {
	Create(ctx context.Context, in *Entity) (err error)
	Update(ctx context.Context, in *Entity) (err error)
	Delete(ctx context.Context, in *Entity) (err error)
}

type QueryResolver func() IQuery
type CommandResolver func() ICommand
