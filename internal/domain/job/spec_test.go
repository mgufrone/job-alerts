package job

import (
	"github.com/stretchr/testify/assert"
	criteria2 "mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"testing"
)

func TestSourceSpec(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         criteria2.ICondition
		shouldError bool
	}{
		{"", nil, true},
		{"tia", criteria2.NewCondition("source", criteria2.Eq, "tia"), true},
		{"something", criteria2.NewCondition("source", criteria2.Eq, "something"), true},
		{"123", criteria2.NewCondition("source", criteria2.Eq, "123"), true},
	}
	for _, c := range testCases {
		r := WhereSource(c.in)
		assert.Equal(t, c.out, r)
	}
}
func TestURLSpec(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         criteria2.ICondition
		shouldError bool
	}{
		{"", nil, true},
		{"tia", criteria2.NewCondition("job_url", criteria2.Eq, "tia"), true},
		{"something", criteria2.NewCondition("job_url", criteria2.Eq, "something"), true},
		{"123", criteria2.NewCondition("job_url", criteria2.Eq, "123"), true},
	}
	for _, c := range testCases {
		r := WhereURL(c.in)
		assert.Equal(t, c.out, r)
	}
}
