package publisher

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Name() string {
	return m.Called().String(0)
}

func (m *MockPublisher) Publish(ctx context.Context, entity *notification.Entity) error {
	return m.Called(ctx, entity).Error(0)
}

func (m *MockPublisher) SetReceiver(entity *user_channel.Entity) {
	m.Called(entity)
}
