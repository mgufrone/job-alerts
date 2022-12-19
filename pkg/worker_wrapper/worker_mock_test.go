package wrapper

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/test_data"
)

type workerMock struct {
	mock.Mock
}

func (w *workerMock) Fetch(ctx context.Context) ([]*job.Entity, error) {
	args := w.Called(ctx)
	return args.Get(0).([]*job.Entity), args.Error(1)
}

func (w *workerMock) FetchJob(ctx context.Context, job2 *job.Entity) (*job.Entity, error) {
	args := w.Called(ctx, job2)
	return args.Get(0).(*job.Entity), args.Error(1)
}

func (w *workerMock) mockFetch(fail bool, msg string, numRes int) {
	var err error
	jbs := make([]*job.Entity, numRes)
	if numRes > 0 {
		for i := 0; i < numRes; i++ {
			jbs[i] = test_data.ValidJob()
		}
	}
	if fail {
		err = errors.New(msg)
	}
	w.On("Fetch", mock.Anything).Return(jbs, err).Once()
}
func (w *workerMock) mockFetchJob(fail bool, msg string) {
	var (
		err error
		jb  *job.Entity
	)
	if fail {
		err = errors.New(msg)
	} else {
		jb = test_data.ValidJob()
	}
	w.On("Fetch", mock.Anything).Return(jb, err).Once()
}
