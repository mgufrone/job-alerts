package notification

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/test_data"
	"testing"
)

type testCaseSubscribe struct {
	in       *channel.Entity
	before   func(mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository)
	after    func(t *testing.T, mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository)
	wantsErr bool
}
type mockValueSubscribe struct {
	mQuery  *channel.MockQueryRepository
	mCmd    *channel.MockCommandRepository
	useCase *UseCase
}

func setupMockSubscribe() *mockValueSubscribe {
	mv := &mockValueSubscribe{
		mQuery: &channel.MockQueryRepository{},
		mCmd:   &channel.MockCommandRepository{},
	}
	mv.useCase = &UseCase{
		chQuery: func() channel.IQueryRepository {
			return mv.mQuery
		},
		chCommand: func() channel.ICommandRepository {
			return mv.mCmd
		},
	}
	return mv
}

func TestNotification_Subscribe(t *testing.T) {
	ch := test_data.ValidChannel(true)
	cases := []testCaseSubscribe{
		{
			nil,
			func(mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository) {

			},
			nil,
			true,
		},
		// test exists
		{
			ch,
			func(mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				mQuery.On("CriteriaBuilder").Return(mCriteria)
				mQuery.On("Apply", mock.Anything).Return(mCriteria)
				mQuery.On("Count", mock.Anything).Return(int64(1), nil)
			},
			func(t *testing.T, mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository) {
				mQuery.AssertNumberOfCalls(t, "Count", 1)

			},
			true,
		},
		// test create entries
		{
			ch,
			func(mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				mQuery.On("CriteriaBuilder").Return(mCriteria)
				mQuery.On("Apply", mock.Anything).Return(mCriteria)
				mQuery.On("Count", mock.Anything).Return(int64(0), nil)

				mCmd.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mQuery *channel.MockQueryRepository, mCmd *channel.MockCommandRepository) {
				mQuery.AssertNumberOfCalls(t, "Count", 1)
				mCmd.AssertNumberOfCalls(t, "Create", 1)
			},
			false,
		},
	}
	for _, c := range cases {
		mck := setupMockSubscribe()
		c.before(mck.mQuery, mck.mCmd)
		err := mck.useCase.Subscribe(context.Background(), c.in)
		if c.wantsErr {
			assert.NotNil(t, err, c)
			continue
		}
		assert.Nil(t, err)
		c.after(t, mck.mQuery, mck.mCmd)
	}
}
