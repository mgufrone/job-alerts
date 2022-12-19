package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUUID(t *testing.T) {
	t.Parallel()
	okID := uuid.NewString()
	testCases := []struct {
		in         []byte
		keys       []string
		shouldFail bool
		out        uuid.UUID
	}{
		{[]byte("null"), []string{}, true, uuid.Nil},
		{[]byte("{}"), []string{}, true, uuid.Nil},
		{[]byte("1640412101"), []string{}, true, uuid.Nil},
		{[]byte("{time:1640412101}"), []string{"time"}, true, uuid.Nil},
		{[]byte(`{"id":null}`), []string{"time"}, true, uuid.Nil},
		{[]byte(`{"id":null}`), []string{"id"}, false, uuid.Nil},
		{[]byte(`{"id":"0000-0000-0000-0000"}`), []string{"id"}, true, uuid.Nil},
		{[]byte(fmt.Sprintf(`{"id":"%s"}`, okID)), []string{"id"}, false, uuid.MustParse(okID)},
	}
	for _, c := range testCases {
		res, err := GetUUID(c.in, c.keys...)
		if c.shouldFail {
			require.NotNil(t, err)
			require.Equal(t, c.out, res)
			continue
		}
		require.NotNil(t, res)
		require.Nil(t, err)
		require.Equal(t, c.out, res)
	}
}
