package user

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEntity_ID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in         uuid.UUID
		shouldFail bool
	}{
		{uuid.Nil, false},
		{uuid.New(), false},
	}
	for _, c := range testCases {
		var u Entity
		err := u.SetID(c.in)
		if c.shouldFail {
			require.NotNil(t, err)
			continue
		}
		require.Nil(t, err)
		require.Equal(t, c.in, u.ID())
	}
}

func TestEntity_AuthID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in         string
		shouldFail bool
	}{
		{"", true},
		{"ok", false},
	}
	for _, c := range testCases {
		var u Entity
		err := u.SetAuthID(c.in)
		if c.shouldFail {
			require.NotNil(t, err)
			continue
		}
		require.Nil(t, err)
		require.Equal(t, c.in, u.AuthID())
	}
}

func TestEntity_Status(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in         Status
		shouldFail bool
	}{
		{0, true},
		{100, true},
		{1, false},
		{Active, false},
		{Deactivated, false},
	}
	for _, c := range testCases {
		var u Entity
		err := u.SetStatus(c.in)
		if c.shouldFail {
			require.NotNil(t, err, c.in)
			continue
		}
		require.Nil(t, err, "input", c)
		require.Equal(t, c.in, u.Status())
	}
}
func TestEntity_Roles(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in         Role
		shouldFail bool
	}{
		{0, true},
		{User, false},
		{Guest | User, false},
		{Admin, false},
		{Admin | User, false},
	}
	for _, c := range testCases {
		var u Entity
		err := u.SetRoles(c.in)
		if c.shouldFail {
			require.NotNil(t, err, c.in)
			continue
		}
		require.Nil(t, err)
		require.Equal(t, c.in, u.Roles())
	}
}
func TestEntity_HasRole(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in    Role
		check Role
		out   bool
	}{
		{0, 0, false},
		{Guest, User, false},
		{Guest | User, Admin, false},
		{Admin | User, User, true},
		{Guest | User, Guest, true},
	}
	for _, c := range testCases {
		var u Entity
		_ = u.SetRoles(c.in)
		require.Equal(t, c.out, u.HasRole(c.check), c.in)
	}
}
