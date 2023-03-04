package user_channel

import (
	"context"
	"github.com/google/uuid"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type IQueryRepository interface {
	// CriteriaBuilder will handle the filter and pagination process
	CriteriaBuilder() criteria.ICriteriaBuilder
	// Apply will submit criteria builder based on
	Apply(cb criteria.ICriteriaBuilder) IQueryRepository
	// GetAll users record
	GetAll(ctx context.Context) (out []*Entity, err error)
	// Count user based on the criteria
	Count(ctx context.Context) (total int64, err error)
	// GetAndCount is shorthand for GetAll and Count
	GetAndCount(ctx context.Context) (out []*Entity, total int64, err error)
	// FindByID filter user by id
	FindByID(ctx context.Context, id uuid.UUID) (out *Entity, err error)
}
type ICommandRepository interface {
	// Create a new user
	Create(ctx context.Context, in *Entity) (err error)
	// Update existing user
	Update(ctx context.Context, in *Entity) (err error)
	// Delete existing user
	Delete(ctx context.Context, in *Entity) (err error)
}
type QueryResolver func() IQueryRepository
type CommandResolver func() ICommandRepository
