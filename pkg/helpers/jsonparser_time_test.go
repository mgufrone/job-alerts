package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTime(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in         []byte
		keys       []string
		shouldFail bool
		out        time.Time
	}{
		{[]byte("null"), []string{}, true, time.Time{}},
		{[]byte("{}"), []string{}, true, time.Time{}},
		{[]byte("1640412101"), []string{}, true, time.Time{}},
		{[]byte("{time:1640412101}"), []string{"time"}, true, time.Time{}},
		{[]byte(`{"time":1640412101000}`), []string{"time"}, false, time.Unix(1640412101, 0)},
	}
	for _, c := range testCases {
		res, err := GetTime(c.in, c.keys...)
		if c.shouldFail {
			assert.NotNil(t, err)
			assert.True(t, res.IsZero())
			continue
		}
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.Equal(t, c.out, res)
	}
}
