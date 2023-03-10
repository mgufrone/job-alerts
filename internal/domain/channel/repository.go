package channel

import (
	"context"
	"github.com/google/uuid"
	criteria2 "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type IQueryRepository interface {
	CriteriaBuilder() criteria2.ICriteriaBuilder
	Apply(cb criteria2.ICriteriaBuilder) IQueryRepository
	// expose-able
	GetAll(ctx context.Context) (out []*Entity, err error)
	Count(ctx context.Context) (total int64, err error)
	GetAndCount(ctx context.Context) (out []*Entity, total int64, err error)
	FindByID(ctx context.Context, id uuid.UUID) (out *Entity, err error)
}

type ICommandRepository interface {
	Create(ctx context.Context, in *Entity) (err error)
	Update(ctx context.Context, in *Entity) (err error)
	Delete(ctx context.Context, in *Entity) (err error)
}

type QueryResolver func() IQueryRepository
type CommandResolver func() ICommandRepository
