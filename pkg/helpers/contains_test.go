package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContains(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		haystack []string
		needle   string
		ok       bool
	}{
		{[]string{}, "", false},
		{[]string{"ok"}, "", false},
		{[]string{"ok"}, "ko", false},
		{[]string{"ok"}, "k", false},
		{[]string{}, "k", false},
		{[]string{"ok", "b", "c", "d", "k"}, "k", true},
	}
	for _, c := range testCases {
		res := Contains(c.haystack, c.needle)
		assert.Equal(t, c.ok, res)
	}
}
