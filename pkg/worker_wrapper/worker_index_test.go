package wrapper

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"testing"
)

func TestWorkerIndex_Run(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		fail      bool
		msg       string
		limit     int
		queryCall bool
		queryFail bool
		cmdCall   bool
	}{
		{true, "something wrong", 0, false, false, false},
		{false, "something wrong", 0, false, false, false},
		{false, "something wrong", 5, true, false, true},
	}
	for _, tc := range testCases {
		var (
			wm    workerMock
			query job.MockQueryRepository
			cmd   job.MockCommandRepository
			cb    criteria.MockCriteria
		)
		wm.mockFetch(tc.fail, tc.msg, tc.limit)
		if tc.cmdCall {
			cmd.On("Create", mock.Anything, mock.Anything).Return(nil)
		}
		cb.On("Where", mock.Anything).Return(&cb)
		cb.On("Where", mock.Anything).Return(&cb)
		cb.On("Or", mock.Anything).Return(&cb)
		cb.On("And", mock.Anything).Return(&cb)
		cb.On("Paginate", mock.Anything, mock.Anything).Return(&cb)
		query.On("CriteriaBuilder").Return(&cb)
		var err error
		if tc.queryFail {
			err = errors.New("something wrong")
		}
		query.On("GetAll", mock.Anything).Return([]*job.Entity{}, err)
		wrk := NewWorkerIndex(&wm, &cmd, &query)
		err1 := wrk.Run(context.Background())
		wm.AssertExpectations(t)
		if tc.fail {
			assert.NotNil(t, err1)
			continue
		}
		if tc.limit == 0 {
			continue
		}
		query.AssertExpectations(t)
		if tc.queryFail {
			assert.NotNil(t, err1)
			continue
		}
		cmd.AssertExpectations(t)
		cb.AssertExpectations(t)
	}
}
