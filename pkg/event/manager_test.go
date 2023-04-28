package event

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"mgufrone.dev/job-alerts/pkg/payload"
	"testing"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Name() string {
	return m.Called().String(0)
}

func (m *mockHandler) Run(ctx context.Context, py payload.Payload) error {
	return m.Called(ctx, py).Error(0)
}

func (m *mockHandler) When(py payload.Payload) bool {
	return m.Called(py).Bool(0)
}

func (m *mockHandler) SubscribeAt() []string {
	call := m.Called()
	return call.Get(0).([]string)
}

func TestSimpleManager_Subscribe(t *testing.T) {
	sm := NewSimpleManager()
	mck := &mockHandler{}
	mck.On("Name").Return("somename")
	mck.On("SubscribeAt").Return([]string{"something"})
	mck.On("Run").Return(nil)
	mck.On("When").Return(true)
	sm.Subscribe("something", mck)
	assert.Len(t, sm.runners, 1)
	assert.Len(t, sm.handlers, 1)
	assert.Equal(t, sm.handlers["something"], []int{0})
	assert.Equal(t, sm.runners[0], mck)
}

func TestSimpleManager_Unsubscribe(t *testing.T) {
	sm := NewSimpleManager()
	mck := &mockHandler{}
	mck.On("Name").Return("somename")
	mck.On("SubscribeAt").Return([]string{"something"})
	mck.On("Run").Return(nil)
	mck.On("When").Return(true)
	idx := sm.Subscribe("something", mck)
	assert.Len(t, sm.runners, 1)
	assert.Len(t, sm.handlers, 1)
	assert.Equal(t, sm.handlers["something"], []int{0})
	assert.Equal(t, sm.runners[0], mck)
	err := sm.Unsubscribe(idx)
	require.Nil(t, err)
	err2 := sm.Unsubscribe(2)
	require.NotNil(t, err2)
	require.Len(t, sm.runners, 1)
	require.Len(t, sm.handlers, 1)
	require.Len(t, sm.handlers["something"], 0)
	require.Equal(t, sm.runners[0], nil)
}

func TestSimpleManager_Publish(t *testing.T) {
	sm := NewSimpleManager()
	mck := &mockHandler{}
	mck.On("Name").Return("somename")
	mck.On("SubscribeAt").Return([]string{"something"})
	mck.On("Run", mock.Anything, mock.Anything).Return(nil)
	mck.On("When", mock.Anything).Return(true)
	mck2 := &mockHandler{}
	mck2.On("Name").Return("somename")
	mck2.On("SubscribeAt").Return([]string{"something"})
	mck2.On("Run", mock.Anything, mock.Anything).Return(nil)
	mck2.On("When", mock.Anything).Return(false)
	mck3 := &mockHandler{}
	mck3.On("Name").Return("somename")
	mck3.On("SubscribeAt").Return([]string{"notregistered"})
	mck3.On("Run", mock.Anything, mock.Anything).Return(nil)
	mck3.On("When", mock.Anything).Return(true)
	sm.Subscribe("something", mck)
	sm.Subscribe("something", mck2)
	sm.Subscribe("notregistered", mck3)
	sm.Publish(context.TODO(), "something", payload.New([]byte(`{}`)))
	mck.AssertNumberOfCalls(t, "When", 1)
	mck2.AssertNumberOfCalls(t, "When", 1)
	mck3.AssertNumberOfCalls(t, "When", 0)
	mck.AssertNumberOfCalls(t, "Run", 1)
	mck2.AssertNumberOfCalls(t, "Run", 0)
	mck3.AssertNumberOfCalls(t, "Run", 0)
}
