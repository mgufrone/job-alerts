package notification

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MutateMock struct {
	mock.Mock
}

func (m *MutateMock) Create(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}

func (m *MutateMock) Update(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}

func (m *MutateMock) Delete(ctx context.Context, in *Entity) (err error) {
	return m.Called(ctx, in).Error(0)
}
