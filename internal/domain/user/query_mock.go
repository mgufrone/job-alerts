package user

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type MockQueryRepository struct {
	mock.Mock
}

func (q *MockQueryRepository) CriteriaBuilder() criteria.ICriteriaBuilder {
	return q.Called().Get(0).(criteria.ICriteriaBuilder)
}

func (q *MockQueryRepository) Apply(cb criteria.ICriteriaBuilder) IQueryRepository {
	q.Called(cb)
	return q
}

func (q *MockQueryRepository) GetAll(ctx context.Context) (out []*Entity, err error) {
	args := q.Called(ctx)
	return args.Get(0).([]*Entity), args.Error(1)
}

func (q *MockQueryRepository) Count(ctx context.Context) (total int64, err error) {
	args := q.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (q *MockQueryRepository) GetAndCount(ctx context.Context) (out []*Entity, total int64, err error) {
	args := q.Called(ctx)
	return args.Get(0).([]*Entity), args.Get(1).(int64), args.Error(2)
}

func (q *MockQueryRepository) FindByID(ctx context.Context, id string) (out *Entity, err error) {
	args := q.Called(ctx, id)
	return args.Get(0).(*Entity), args.Get(1).(error)
}
