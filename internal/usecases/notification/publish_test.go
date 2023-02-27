package notification

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
	"mgufrone.dev/job-alerts/internal/services/publisher"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/test_data"
	"testing"
)

type testCasePublish struct {
	in       *notification.Entity
	before   func(mp *mockValuePublish, ent *notification.Entity)
	after    func(t *testing.T, mp *mockValuePublish, ent *notification.Entity)
	wantsErr bool
}
type mockValuePublish struct {
	chQuery     *channel.MockQueryRepository
	chCmd       *channel.MockCommandRepository
	nQuery      *notification.MockQueryRepository
	nCmd        *notification.MockCommandRepository
	jQuery      *job.MockQueryRepository
	usrQuery    *user.MockQueryRepository
	usrChQuery  *user_channel.MockQueryRepository
	collections []*publisher.MockPublisher
	useCase     *UseCase
}

func setupMockPublish() *mockValuePublish {
	mv := &mockValuePublish{
		chQuery:    channel.NewMockQueryRepository(),
		chCmd:      channel.NewMockCommandRepository(),
		nQuery:     &notification.MockQueryRepository{},
		nCmd:       &notification.MockCommandRepository{},
		jQuery:     &job.MockQueryRepository{},
		usrQuery:   &user.MockQueryRepository{},
		usrChQuery: &user_channel.MockQueryRepository{},
	}
	mv.useCase = &UseCase{
		chQuery: func() channel.IQueryRepository {
			return mv.chQuery
		},
		chCommand: func() channel.ICommandRepository {
			return mv.chCmd
		},
		notifQuery: func() notification.IQueryRepository {
			return mv.nQuery
		},
		notifCmd: func() notification.ICommandRepository {
			return mv.nCmd
		},
		usrChQuery: func() user_channel.IQueryRepository {
			return mv.usrChQuery
		},
		usrQuery: func() user.IQueryRepository {
			return mv.usrQuery
		},
		jQuery: func() job.IQueryRepository {
			return mv.jQuery
		},
	}
	return mv
}
func TestNotification_Publish(t *testing.T) {
	t.Parallel()
	// this method should be called when the notification is created
	// first, it will check user and job is actually exists (in parallel)
	// if the job, channel or the user is not found, delete the notification
	// if user is active, it will need to retrieve the user notification channels based on the channel subscription (email, in-app, telegram)
	usr := &user.Entity{}
	usr.SetID(uuid.New())
	ch, jb := test_data.ValidChannel(true), test_data.ValidJob()
	cases := []testCasePublish{
		{
			nil,
			func(mp *mockValuePublish, entity *notification.Entity) {

			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.chQuery.AssertNumberOfCalls(t, "Count", 0)
				mp.nQuery.AssertNumberOfCalls(t, "Count", 0)
				mp.jQuery.AssertNumberOfCalls(t, "Count", 0)
			},
			true,
		},
		{
			test_data.ValidNotification(usr, jb, ch),
			func(mp *mockValuePublish, entity *notification.Entity) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				//
				mCriteria1 := &criteria.MockCriteria{}
				mCriteria1.On("Or", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("Where", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("And", mock.Anything, mock.Anything).Return(mCriteria1)
				mp.chQuery.On("CriteriaBuilder").Return(mCriteria1)
				mp.chQuery.On("Apply", mock.Anything).Return(mCriteria1)
				mp.chQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.jQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.jQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.jQuery.On("Count", mock.Anything).Return(int64(0), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.nCmd.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.jQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.chQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.nCmd.AssertNumberOfCalls(t, "Delete", 1)
			},
			true,
		},
		{
			test_data.ValidNotification(usr, jb, ch),
			func(mp *mockValuePublish, entity *notification.Entity) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				//
				mCriteria1 := &criteria.MockCriteria{}
				mCriteria1.On("Or", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("Where", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("And", mock.Anything, mock.Anything).Return(mCriteria1)
				mp.chQuery.On("CriteriaBuilder").Return(mCriteria1)
				mp.chQuery.On("Apply", mock.Anything).Return(mCriteria1)
				mp.chQuery.On("Count", mock.Anything).Return(int64(0), nil)
				//
				mp.jQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.jQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.jQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)

				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.nCmd.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.jQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.chQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.nCmd.AssertNumberOfCalls(t, "Delete", 1)
			},
			true,
		},
		{
			test_data.ValidNotification(usr, jb, ch),
			func(mp *mockValuePublish, entity *notification.Entity) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				//
				mCriteria1 := &criteria.MockCriteria{}
				mCriteria1.On("Or", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("Where", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("And", mock.Anything, mock.Anything).Return(mCriteria1)
				mp.chQuery.On("CriteriaBuilder").Return(mCriteria1)
				mp.chQuery.On("Apply", mock.Anything).Return(mCriteria1)
				mp.chQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.jQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.jQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.jQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(0), nil)
				//
				mp.nCmd.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.jQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.chQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.nCmd.AssertNumberOfCalls(t, "Delete", 1)
				mp.nCmd.AssertCalled(t, "Delete", mock.Anything, entity)
			},
			true,
		},
		{
			test_data.ValidNotification(usr, jb, ch),
			func(mp *mockValuePublish, entity *notification.Entity) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				//
				mCriteria1 := &criteria.MockCriteria{}
				mCriteria1.On("Or", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("Where", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("And", mock.Anything, mock.Anything).Return(mCriteria1)
				mp.chQuery.On("CriteriaBuilder").Return(mCriteria1)
				mp.chQuery.On("Apply", mock.Anything).Return(mCriteria1)
				mp.chQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.jQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.jQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.jQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("GetAndCount", mock.Anything).Return(int64(1), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.usrChQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrChQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrChQuery.On("GetAndCount", mock.Anything).Return(nil, int64(0), nil)
				//
			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.jQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.chQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.nCmd.AssertNumberOfCalls(t, "Delete", 0)
				mp.usrChQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
			},
			true,
		},
		{
			test_data.ValidNotification(usr, jb, ch),
			func(mp *mockValuePublish, entity *notification.Entity) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)
				//
				mCriteria1 := &criteria.MockCriteria{}
				mCriteria1.On("Or", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("Where", mock.Anything, mock.Anything).Return(mCriteria1)
				mCriteria1.On("And", mock.Anything, mock.Anything).Return(mCriteria1)
				mp.chQuery.On("CriteriaBuilder").Return(mCriteria1)
				mp.chQuery.On("Apply", mock.Anything).Return(mCriteria1)
				mp.chQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.jQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.jQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.jQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				mp.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)
				//
				usrCh := test_data.ValidUserChannel(entity.User(), entity.Channel().Channels()[0])
				mp.usrChQuery.On("CriteriaBuilder").Return(mCriteria)
				mp.usrChQuery.On("Apply", mock.Anything).Return(mCriteria)
				mp.usrChQuery.On("GetAndCount", mock.Anything).Return([]*user_channel.Entity{usrCh}, int64(1), nil)
				//
				pb := &publisher.MockPublisher{}
				pb.On("Name").Return(usrCh.ChannelType())
				pb.On("SetReceiver", mock.Anything).Return()
				pb.On("Publish", mock.Anything, mock.Anything).Return(nil)
				mp.collections = append(mp.collections, pb)
				mp.useCase.publishers = publisher.New([]publisher.Publisher{pb})
			},
			func(t *testing.T, mp *mockValuePublish, entity *notification.Entity) {
				mp.jQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.chQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mp.nCmd.AssertNumberOfCalls(t, "Delete", 0)
				mp.usrChQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				mp.collections[0].AssertNumberOfCalls(t, "Publish", 1)
				mp.collections[0].AssertExpectations(t)
			},
			false,
		},
	}
	for _, c := range cases {
		mp := setupMockPublish()
		ctx := context.Background()
		c.before(mp, c.in)
		err := mp.useCase.Publish(ctx, c.in)
		c.after(t, mp, c.in)
		if c.wantsErr {
			assert.Error(t, err)
		}
	}
}
