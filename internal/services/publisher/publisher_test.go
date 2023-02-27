package publisher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockPublisher(name string) *MockPublisher {
	m := &MockPublisher{}
	m.On("Name").Return(name)
	return m
}
func TestCollection_Get(t *testing.T) {
	mp := mockPublisher("world")
	cases := []struct {
		in         []Publisher
		getName    string
		wants      *MockPublisher
		wantsErr   bool
		wantsCount int
	}{
		{
			nil,
			"",
			nil,
			true,
			0,
		},
		{
			[]Publisher{mockPublisher("hello")},
			"world",
			nil,
			true,
			1,
		},
		{
			[]Publisher{mockPublisher("hello"), mp},
			"world",
			mp,
			false,
			2,
		},
		{
			[]Publisher{mockPublisher("hello"), mp, mockPublisher("world")},
			"world",
			mp,
			false,
			2,
		},
	}

	for _, c := range cases {
		coll := New(c.in)
		pub, err := coll.Get(c.getName)
		assert.Equal(t, c.wantsCount, len(coll.publishers))
		if c.wantsErr {
			assert.Nil(t, pub)
			assert.NotNil(t, err)
			continue
		}
		assert.NotNil(t, pub)
		assert.Nil(t, err)
		assert.Equal(t, c.getName, pub.Name())
		assert.Equal(t, c.wants, pub)
	}
}
