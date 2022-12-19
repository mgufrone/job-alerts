package channel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntity_ScheduleAt(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in          string
		out         string
		shouldError bool
	}{
		{"", "", true},
		{"😁", "", true},
		{"*", "", true},
		{"* * * *", "", true},
		{"* * * *", "", true},
		{"* * * * 8", "", true},
		{"60 * * * 8", "", true},
		{"60 25 * * 8", "", true},
		{"60 -1 * * 8", "", true},
		{"60 -1 32 * 8", "", true},
		{"60 -1 -1 * 8", "", true},
		{"60 -1 -1 0 8", "", true},
		{"60 -1 -1 13 8", "", true},
		{"59 1 * * *", "59 1 * * *", false},
		{"59 1 * * *", "59 1 * * *", false},
		{"0/5 1 * * *", "0/5 1 * * *", false},
		{"*/5 1 * * *", "*/5 1 * * *", false},
		{"5,10,15 * * * *", "5,10,15 * * * *", false},
		{"5,10,15 */2 * * *", "5,10,15 */2 * * *", false},
	}
	e := Entity{}
	for _, c := range testCases {
		err := e.SetScheduleAt(c.in)
		if c.shouldError {
			assert.NotNil(t, err, c.in)
			continue
		}
		assert.Equal(t, c.out, e.ScheduleAt())
	}
}
