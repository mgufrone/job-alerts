package job

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/logger"
	"testing"
)

func TestUseCase_List(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	cb := &criteria.MockCriteria{}
	app := fx.New(
		logger.Module,
		job.MockModule,
		notification.MockModule,
		fx.Provide(
			NewUseCase,
		),
		// do mock preps here
		fx.Invoke(func(repository *job.MockQueryRepository) {
			cb.On("Where", mock.Anything).Return(cb)
			cb.On("Or", mock.Anything).Return(cb)
			cb.On("And", mock.Anything).Return(cb)
			cb.On("Paginate", mock.Anything, mock.Anything).Return(cb)
			cb.On("Order", mock.Anything, mock.Anything).Return(cb)

			repository.On("CriteriaBuilder").Return(cb)
			repository.On("And", mock.Anything, mock.Anything).Return(repository)
			repository.On("Apply", mock.Anything).Return(repository)
			repository.On("GetAndCount", mock.Anything).Return([]*job.Entity{}, int64(0), nil)
		}),
		// do test here
		fx.Invoke(func(uc *UseCase) {
			input := ListInput{
				Keyword: "devops",
			}
			jbs, total, err := uc.List(ctx, input)
			assert.Nil(t, err)
			assert.Equal(t, int64(0), total)
			assert.Len(t, jbs, 0)
		}),
		// assertion on mocks
		fx.Invoke(func(jQuery *job.MockQueryRepository) {
			jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
			cb.AssertExpectations(t)
		}),
	)
	app.Start(ctx)
	app.Stop(ctx)
}

func TestUseCase_List2(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	cb := &criteria.MockCriteria{}
	app := fx.New(
		logger.Module,
		job.MockModule,
		notification.MockModule,
		fx.Provide(
			NewUseCase,
		),
		// do mock preps here
		fx.Invoke(func(repository *job.MockQueryRepository) {
			cb.On("Where", mock.Anything).Return(cb)
			cb.On("Or", mock.Anything).Return(cb)
			cb.On("And", mock.Anything, mock.Anything).Return(cb)
			cb.On("Paginate", mock.Anything, mock.Anything).Return(cb)
			cb.On("Order", mock.Anything, mock.Anything).Return(cb)

			repository.On("CriteriaBuilder").Return(cb)
			repository.On("And", mock.Anything, mock.Anything).Return(repository)
			repository.On("Apply", mock.Anything).Return(repository)
			repository.On("GetAndCount", mock.Anything).Return([]*job.Entity{}, int64(0), nil)
		}),
		// do test here
		fx.Invoke(func(uc *UseCase) {
			input := ListInput{
				Keyword: "devops",
				Skills:  []string{"gcp", "aws"},
			}
			jbs, total, err := uc.List(ctx, input)
			assert.Nil(t, err)
			assert.Equal(t, int64(0), total)
			assert.Len(t, jbs, 0)
		}),
		// assertion on mocks
		fx.Invoke(func(jQuery *job.MockQueryRepository) {
			jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
			cb.AssertExpectations(t)
		}),
	)
	app.Start(ctx)
	app.Stop(ctx)
}

// TestUseCase_List3 will verify the jobs query from notifications
func TestUseCase_List3(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	cb := &criteria.MockCriteria{}
	app := fx.New(
		logger.Module,
		job.MockModule,
		notification.MockModule,
		fx.Provide(
			NewUseCase,
		),
		// do mock preps here
		fx.Invoke(func(repository *job.MockQueryRepository, notifQuery *notification.MockQueryRepository) {
			cb.On("Where", mock.Anything).Return(cb)
			cb.On("Or", mock.Anything).Return(cb)
			cb.On("And", mock.Anything, mock.Anything).Return(cb)
			cb.On("Paginate", mock.Anything, mock.Anything).Return(cb)
			cb.On("Order", mock.Anything, mock.Anything).Return(cb)
			var (
				jb1, jb2       job.Entity
				notif1, notif2 notification.Entity
			)
			jb1.SetID(uuid.New())
			jb2.SetID(uuid.New())
			notif1.SetJob(&jb1)
			notif2.SetJob(&jb2)
			notifQuery.On("CriteriaBuilder").Return(cb)
			notifQuery.On("And", mock.Anything, mock.Anything).Return(notifQuery)
			notifQuery.On("Apply", mock.Anything).Return(notifQuery)
			notifQuery.On("GetAll", mock.Anything).Return([]*notification.Entity{&notif1, &notif2}, nil)

			repository.On("CriteriaBuilder").Return(cb)
			repository.On("And", mock.Anything, mock.Anything).Return(repository)
			repository.On("Apply", mock.Anything).Return(repository)
			repository.On("GetAndCount", mock.Anything).Return([]*job.Entity{}, int64(0), nil)
		}),
		// do test here
		fx.Invoke(func(uc *UseCase) {
			var usr user.Entity
			usr.SetID(uuid.New())
			input := ListInput{
				Keyword: "devops",
				User:    &usr,
				Skills:  []string{"gcp", "aws"},
			}
			jbs, total, err := uc.List(ctx, input)
			assert.Nil(t, err)
			assert.Equal(t, int64(0), total)
			assert.Len(t, jbs, 0)
		}),
		// assertion on mocks
		fx.Invoke(func(jQuery *job.MockQueryRepository, nQuery *notification.MockQueryRepository) {
			jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
			nQuery.AssertNumberOfCalls(t, "GetAll", 1)
			cb.AssertExpectations(t)
		}),
	)
	app.Start(ctx)
	app.Stop(ctx)
}
