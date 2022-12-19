package notification

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type QueryMock struct {
	mock.Mock
}

func (q *QueryMock) CriteriaBuilder() criteria.ICriteriaBuilder {
	return q.Called().Get(0).(criteria.ICriteriaBuilder)
}

func (q *QueryMock) Apply(cb criteria.ICriteriaBuilder) IQueryRepository {
	q.Called(cb)
	return q
}

func (q *QueryMock) GetAll(ctx context.Context) (out []*Entity, err error) {
	args := q.Called(ctx)
	return args.Get(0).([]*Entity), args.Error(1)
}

func (q *QueryMock) Count(ctx context.Context) (total int64, err error) {
	args := q.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (q *QueryMock) GetAndCount(ctx context.Context) (out []*Entity, total int64, err error) {
	args := q.Called(ctx)
	return args.Get(0).([]*Entity), args.Get(1).(int64), args.Error(2)
}

func (q *QueryMock) FindByID(ctx context.Context, id string) (out *Entity, err error) {
	args := q.Called(ctx, id)
	return args.Get(0).(*Entity), args.Get(1).(error)
}
