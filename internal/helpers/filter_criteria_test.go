package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/internal/domain/channel"
	"mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"testing"
)

/*
*
- get jobs with the criteria set in the found channels, supported filter:
  - recent jobs (prior to 1 month from the created date. beyond date will be presumed as expired/closed)
  - keyword (against title, role, companyName, or description)
  - tags/skills
  - source (upwork, weworkremotely)
  - range salary
  - is remote
  - job/contract type (full time, part time, contract/freelance)
*/
func TestFilterCriteria(t *testing.T) {
	cases := []struct {
		in    []byte
		after func(mc *criteria.MockCriteria)
	}{
		{
			[]byte(`[]`),
			func(mc *criteria.MockCriteria) {
				mc.AssertNumberOfCalls(t, "Where", 0)
			},
		},
		{
			[]byte(`null`),
			func(mc *criteria.MockCriteria) {
				mc.AssertNumberOfCalls(t, "Where", 0)
			},
		},
		{
			[]byte(`{"skills":"devops,aws"}`),
			func(mc *criteria.MockCriteria) {
				mc.AssertNumberOfCalls(t, "And", 1)
				mc.AssertNumberOfCalls(t, "Or", 0)
				mc.AssertCalled(t, "Where",
					job.WhereTags([]string{"devops", "aws"}),
				)
			},
		},
		{
			[]byte(`{"keyword":"senior devops"}`),
			func(mc *criteria.MockCriteria) {
				mc.AssertNumberOfCalls(t, "And",
					1,
				)
				mc.AssertNumberOfCalls(t, "Or",
					1,
				)
				for _, call := range mc.Calls {
					if call.Method == "Where" {
						cond, _ := call.Arguments[0].(criteria.ICondition)
						if cond.Field() == "role" {
							assert.Equal(t,
								job.WhereRoleContains("senior devops"),
								cond,
							)
						}
						if cond.Field() == "description" {
							assert.Equal(t,
								job.WhereDescriptionContains("senior devops"),
								cond,
							)
						}
					}
				}
			},
		},
		{
			[]byte(`{"keyword":"senior devops", "skills":"terraform,gcp,aws,azure"}`),
			func(mc *criteria.MockCriteria) {
				mc.AssertNumberOfCalls(t, "And",
					1,
				)
				mc.AssertNumberOfCalls(t, "Or",
					1,
				)
				mc.AssertNumberOfCalls(t, "Where",
					3,
				)
				for _, call := range mc.Calls {
					if call.Method == "And" {
						assert.Len(t, call.Arguments[0], 2)
					}
					if call.Method == "Where" {
						cond, _ := call.Arguments[0].(criteria.ICondition)
						fmt.Println("field name", cond.Field())
						if cond.Field() == "role" {
							assert.Equal(t,
								job.WhereRoleContains("senior devops"),
								cond,
							)
						}
						if cond.Field() == "description" {
							assert.Equal(t,
								job.WhereDescriptionContains("senior devops"),
								cond,
							)
						}
						if cond.Field() == "tags" {
							assert.Equal(t,
								job.WhereTags([]string{"terraform", "gcp", "aws", "azure"}),
								cond,
							)
						}
					}
				}
			},
		},
		//{
		//	[]byte(`{"salary":"5000"}`),
		//	func(mc *criteria.MockCriteria) {
		//		mc.AssertNumberOfCalls(t, "And", 1)
		//		mc.AssertNumberOfCalls(t, "And", 0)
		//		mc.AssertCalled(t, "Where",
		//			job.WhereTags([]string{"devops", "aws"}),
		//		)
		//	},
		//},
	}
	for _, c := range cases {
		var (
			ch channel.Entity
			cb = &criteria.MockCriteria{}
		)
		cb.On("Where", mock.Anything).Return(cb)
		cb.On("Or", mock.Anything).Return(cb)
		cb.On("And", mock.Anything).Return(cb)
		cb.On("And", mock.Anything, mock.Anything).Return(cb)
		ch.SetCriterias(c.in)
		cr := FilterToCriteria(&ch, cb, cb)
		c.after(cr.(*criteria.MockCriteria))
	}
}
