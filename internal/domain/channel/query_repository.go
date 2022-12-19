package channel

import (
	"context"
	"github.com/stretchr/testify/mock"
	criteria2 "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
)

type MockQueryRepository struct {
	mock.Mock
}

func (m *MockQueryRepository) FindByID(ctx context.Context, id string) (out *Entity, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Entity), args.Error(1)
}

func (m *MockQueryRepository) CriteriaBuilder() criteria2.ICriteriaBuilder {
	return m.Called().Get(0).(criteria2.ICriteriaBuilder)
}

func (m *MockQueryRepository) Apply(cb criteria2.ICriteriaBuilder) IQueryRepository {
	m.Called(cb)
	return m
}

func (m *MockQueryRepository) GetAll(ctx context.Context) (out []*Entity, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Entity), args.Error(1)
}

func (m *MockQueryRepository) Count(ctx context.Context) (total int64, err error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQueryRepository) GetAndCount(ctx context.Context) (out []*Entity, total int64, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Entity), args.Get(1).(int64), args.Error(2)
}

func NewMockQueryRepository() *MockQueryRepository {
	return &MockQueryRepository{}
}
