package notification

import (
	"context"
	"fmt"
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

type testCaseLoop struct {
	in       *notification.Entity
	before   func(mp *mockValueLoop, ent *notification.Entity)
	after    func(t *testing.T, mp *mockValueLoop, ent *notification.Entity)
	wantsErr bool
}
type mockValueLoop struct {
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

func getMockCall(m *mock.Mock, methodName string) *mock.Call {
	for _, c := range m.Calls {
		if c.Method == methodName {
			return &c
		}
	}
	return nil
}

func setupMockLoop() *mockValueLoop {
	mv := &mockValueLoop{
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
func TestUseCase_Looper(t *testing.T) {
	t.Parallel()
	// this method will be responsible to run through all active channels
	// and query the jobs based on the found channels
	/**
	query channels based on the active state:
	- if channel is active
	- if user is not deactivated
	- if user_channel registered in the notification is active
	for simplicity, it will run query with these criteria:
	- get jobs with the criteria set in the found channels, supported filter:
		- recent jobs (prior to 1 month from the created date. beyond date will be presumed as expired/closed)
		- keyword (against title, role, companyName, or description)
		- tags/skills
		- source (upwork, weworkremotely)
		- range salary
		- is remote
		- job/contract type (full time, part time, contract/freelance)
	- for now, any defined criteria will be glued with AND. so no customization control flow for now.
	*/
	// if match found, check if the job is already listed in the notifications table
	// else, create a new entry but it will also adhere the schedule if set in the channel
	// ideally, the notification repository will have observer that will publish it to the user channels
	cases := []struct {
		before func(mv *mockValueLoop)
		after  func(t *testing.T, mv *mockValueLoop)
	}{
		// case 1: found no active channels
		{
			func(mv *mockValueLoop) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)

				mv.chQuery.On("CriteriaBuilder").Return(mCriteria)
				mv.chQuery.On("Apply", mock.Anything).Return(mv.chQuery)
				mv.chQuery.On("GetAndCount", mock.Anything).Return(nil, int64(0), nil)

				mv.usrQuery.On("CriteriaBuilder").Return(mCriteria)
				mv.usrQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.usrQuery.On("Count", mock.Anything).Return(int64(0), nil)
			},
			func(t *testing.T, mv *mockValueLoop) {
				mv.chQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				calls := mv.chQuery.Calls
				var crit *criteria.MockCriteria
				for _, c := range calls {
					if c.Method == "Apply" {
						crit = c.Arguments[0].(*criteria.MockCriteria)
					}
				}
				if crit != nil {
					crit.AssertNumberOfCalls(t, "Where", 1)
					arg := crit.Calls[0].Arguments[0].(criteria.ICondition)
					assert.Equal(t, arg.Field(), "is_active")
					assert.Equal(t, arg.Operator(), criteria.Eq)
					assert.Equal(t, arg.Value(), true)
				}

				mv.jQuery.AssertNumberOfCalls(t, "Count", 0)
				mv.usrQuery.AssertNumberOfCalls(t, "Count", 0)
				mv.usrChQuery.AssertNumberOfCalls(t, "Count", 0)
				// the main idea would be this one
				mv.nCmd.AssertNumberOfCalls(t, "Create", 0)
			},
		},
		// case 2: found a channel but no jobs active
		{
			func(mv *mockValueLoop) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)

				mv.chQuery.On("CriteriaBuilder").Return(mCriteria)
				mv.chQuery.On("Apply", mock.Anything).Return(mv.chQuery)
				mv.chQuery.On("GetAndCount", mock.Anything).Return([]*channel.Entity{
					test_data.ValidChannel(true),
				}, int64(1), nil)

				mCriteria2 := &criteria.MockCriteria{}
				mCriteria2.On("Or", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("And", mock.Anything, mock.Anything).Return(mCriteria2)

				mv.usrQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.usrQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)

				mv.jQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.jQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.jQuery.On("GetAndCount", mock.Anything).Return([]*job.Entity{}, int64(0), nil)
			},
			func(t *testing.T, mv *mockValueLoop) {
				mv.chQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				calls := mv.chQuery.Calls
				var crit *criteria.MockCriteria
				for _, c := range calls {
					if c.Method == "Apply" {
						crit = c.Arguments[0].(*criteria.MockCriteria)
					}
				}
				if crit != nil {
					crit.AssertNumberOfCalls(t, "Where", 1)
					arg := crit.Calls[0].Arguments[0].(criteria.ICondition)
					assert.Equal(t, arg.Field(), "is_active")
					assert.Equal(t, arg.Operator(), criteria.Eq)
					assert.Equal(t, arg.Value(), true)
				}

				mv.jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				mv.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mv.usrChQuery.AssertNumberOfCalls(t, "Count", 0)
				// the main idea would be this one
				mv.nCmd.AssertNumberOfCalls(t, "Create", 0)
			},
		},
		// case 3: found a channel and active jobs
		{
			func(mv *mockValueLoop) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)

				mv.chQuery.On("CriteriaBuilder").Return(mCriteria)
				mv.chQuery.On("Apply", mock.Anything).Return(mv.chQuery)
				mv.chQuery.On("GetAndCount", mock.Anything).Return([]*channel.Entity{
					test_data.ValidChannel(true),
				}, int64(1), nil)

				mCriteria2 := &criteria.MockCriteria{}
				mCriteria2.On("Or", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("And", mock.Anything, mock.Anything).Return(mCriteria2)

				mv.usrQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.usrQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)

				mv.jQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.jQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.jQuery.On("GetAndCount", mock.Anything).Return([]*job.Entity{test_data.ValidJob(), test_data.ValidJob()}, int64(2), nil)

				mv.nQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.nQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.nQuery.On("Count", mock.Anything).Once().Return(int64(1), nil)
				mv.nQuery.On("Count", mock.Anything).Return(int64(0), nil)

				mv.nCmd.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mv *mockValueLoop) {
				mv.chQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				calls := mv.chQuery.Calls
				var crit *criteria.MockCriteria
				for _, c := range calls {
					if c.Method == "Apply" {
						crit = c.Arguments[0].(*criteria.MockCriteria)
					}
				}
				if crit != nil {
					crit.AssertNumberOfCalls(t, "Where", 1)
					arg := crit.Calls[0].Arguments[0].(criteria.ICondition)
					assert.Equal(t, arg.Field(), "is_active")
					assert.Equal(t, arg.Operator(), criteria.Eq)
					assert.Equal(t, arg.Value(), true)
				}

				mv.jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				for _, c := range mv.jQuery.Calls {
					if c.Method == "Apply" {
						criCalls := c.Arguments[0].(*criteria.MockCriteria).Calls
						for _, ca := range criCalls {
							if ca.Method == "Where" {
								arg, ok := ca.Arguments[0].(criteria.ICondition)
								if ok {
									fmt.Println("passed arg", arg.ToString())
								}
							}
						}
					}
				}
				mv.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mv.usrChQuery.AssertNumberOfCalls(t, "Count", 0)
				mv.nQuery.AssertNumberOfCalls(t, "Count", 2)
				// the main idea would be this one
				mv.nCmd.AssertNumberOfCalls(t, "Create", 1)
			},
		},
		// case 3: found a channel and active jobs
		{
			func(mv *mockValueLoop) {
				mCriteria := &criteria.MockCriteria{}
				mCriteria.On("Or", mock.Anything, mock.Anything).Return(mCriteria)
				mCriteria.On("Where", mock.Anything).Return(mCriteria)
				mCriteria.On("And", mock.Anything, mock.Anything).Return(mCriteria)

				mv.chQuery.On("CriteriaBuilder").Return(mCriteria)
				mv.chQuery.On("Apply", mock.Anything).Return(mv.chQuery)
				mv.chQuery.On("GetAndCount", mock.Anything).Return([]*channel.Entity{
					test_data.ValidChannel(true),
				}, int64(1), nil)

				mCriteria2 := &criteria.MockCriteria{}
				mCriteria2.On("Or", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything).Return(mCriteria2)
				mCriteria2.On("Where", mock.Anything, mock.Anything).Return(mCriteria2)
				mCriteria2.On("And", mock.Anything, mock.Anything).Return(mCriteria2)

				mv.usrQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.usrQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.usrQuery.On("Count", mock.Anything).Return(int64(1), nil)

				mv.jQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.jQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.jQuery.On("GetAndCount", mock.Anything).Return([]*job.Entity{test_data.ValidJob(), test_data.ValidJob(), test_data.ValidJob()}, int64(2), nil)

				mv.nQuery.On("CriteriaBuilder").Return(mCriteria2)
				mv.nQuery.On("Apply", mock.Anything).Return(mv.usrQuery)
				mv.nQuery.On("Count", mock.Anything).Once().Return(int64(1), nil)
				mv.nQuery.On("Count", mock.Anything).Return(int64(0), nil)

				mv.nCmd.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			func(t *testing.T, mv *mockValueLoop) {
				mv.chQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				calls := mv.chQuery.Calls
				var crit *criteria.MockCriteria
				for _, c := range calls {
					if c.Method == "Apply" {
						crit = c.Arguments[0].(*criteria.MockCriteria)
					}
				}
				if crit != nil {
					crit.AssertNumberOfCalls(t, "Where", 1)
					arg := crit.Calls[0].Arguments[0].(criteria.ICondition)
					assert.Equal(t, arg.Field(), "is_active")
					assert.Equal(t, arg.Operator(), criteria.Eq)
					assert.Equal(t, arg.Value(), true)
				}

				mv.jQuery.AssertNumberOfCalls(t, "GetAndCount", 1)
				for _, c := range mv.jQuery.Calls {
					if c.Method == "Apply" {
						criCalls := c.Arguments[0].(*criteria.MockCriteria).Calls
						for _, ca := range criCalls {
							if ca.Method == "Where" {
								arg, ok := ca.Arguments[0].(criteria.ICondition)
								if ok {
									fmt.Println("passed criteria", arg.ToString())
								}
							}
						}
					}
				}
				mv.usrQuery.AssertNumberOfCalls(t, "Count", 1)
				mv.usrChQuery.AssertNumberOfCalls(t, "Count", 0)
				mv.nQuery.AssertNumberOfCalls(t, "Count", 3)
				// the main idea would be this one
				mv.nCmd.AssertNumberOfCalls(t, "Create", 2)
			},
		},
	}
	for _, c := range cases {
		ml := setupMockLoop()
		ctx := context.Background()
		c.before(ml)
		ml.useCase.Loop(ctx)
		c.after(t, ml)
	}
}
