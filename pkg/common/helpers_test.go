package common

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockTry struct {
	mock.Mock
}

func (m *mockTry) Run() error {
	return m.Called().Error(0)
}
func TestTryOrError(t *testing.T) {
	t.Parallel()
	first := new(mockTry)
	first.On("Run").Once().Return(errors.New("error"))
	second := new(mockTry)
	second.On("Run").Once().Return(errors.New("error"))
	err := TryOrError(first.Run, second.Run)
	require.NotNil(t, err)
	first.AssertCalled(t, "Run")
	second.AssertNotCalled(t, "Run")
	first.Calls = []mock.Call{}
	second.Calls = []mock.Call{}
	first.ExpectedCalls = []*mock.Call{}
	second.ExpectedCalls = []*mock.Call{}

	first.On("Run").Once().Return(nil)
	second.On("Run").Once().Return(nil)
	err = TryOrError(first.Run, second.Run)
	first.AssertCalled(t, "Run")
	second.AssertCalled(t, "Run")
	require.Nil(t, err)
}
