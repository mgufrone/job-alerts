package notification

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockCommandRepository struct {
	mock.Mock
}

func (m *MockCommandRepository) Create(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}

func (m *MockCommandRepository) Update(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}

func (m *MockCommandRepository) Delete(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}
