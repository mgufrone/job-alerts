package wrapper

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/job"
	mocks2 "mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/test_data"
	"testing"
)

func TestWorkerDataJob_Run(t *testing.T) {
	t.Parallel()
	validJob := test_data.ValidJob()
	testCases := []struct {
		in        *job.Entity
		failParse bool
		failFetch bool
	}{
		{
			nil, true, false,
		},
		{
			validJob, false, false,
		},
	}
	for _, c := range testCases {
		var (
			cmd   mocks2.MockCommandRepository
			query mocks2.MockQueryRepository
			cb    criteria.MockCriteria
			err2  error
		)
		cb.On("Where", mock.Anything).Return(&cb)
		query.On("CriteriaBuilder").Return(&cb)
		query.On("Apply", mock.Anything).Return(&cb)
		if c.failFetch {
			err2 = errors.New("just an error. not a big deal")
		}
		query.On("Count", mock.Anything).Return(int64(0), err2)
		cmd.On("Create", mock.Anything, mock.Anything).Return(err2)
		hndlr := NewWorkerDataJob(&cmd, &query)
		err := hndlr.Run(context.Background(), c.in)
		if c.failParse {
			assert.NotNil(t, err)
			continue
		}
		query.AssertExpectations(t)
		if c.failFetch {
			assert.NotNil(t, err)
			continue
		}
		cmd.AssertExpectations(t)
	}
}
