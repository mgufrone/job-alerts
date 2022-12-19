package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinute(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{
			"a",
			false,
		},
		{
			"",
			false,
		},
		{
			"-1",
			false,
		},
		{
			"0-60",
			false,
		},
		{
			"30-20",
			false,
		},
		{
			"10,20,30,15",
			false,
		},
		{
			"*",
			true,
		},
		{
			"1",
			true,
		},
		{
			"1,5,9",
			true,
		},
		{
			"0-58",
			true,
		},
		{
			"0-59",
			true,
		},
	}
	for _, c := range cases {
		res := isValidMinute(c.in)
		assert.Equal(t, c.out, res, c.in)
	}
}

func TestHour(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{
			"a",
			false,
		},
		{
			"",
			false,
		},
		{
			"-1",
			false,
		},
		{
			"0-24",
			false,
		},
		{
			"20-15",
			false,
		},
		{
			"10,20,22,15",
			false,
		},
		{
			"*",
			true,
		},
		{
			"1",
			true,
		},
		{
			"1,5,9",
			true,
		},
		{
			"0-22",
			true,
		},
		{
			"0-23",
			true,
		},
	}
	for _, c := range cases {
		res := isValidHour(c.in)
		assert.Equal(t, c.out, res, c.in)
	}
}
func TestDate(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{
			"a",
			false,
		},
		{
			"",
			false,
		},
		{
			"-1",
			false,
		},
		{
			"0-34",
			false,
		},
		{
			"20-15",
			false,
		},
		{
			"10,20,31,15",
			false,
		},
		{
			"*",
			true,
		},
		{
			"1",
			true,
		},
		{
			"1,5,9",
			true,
		},
		{
			"1-22",
			true,
		},
		{
			"2-31",
			true,
		},
	}
	for _, c := range cases {
		res := isValidDate(c.in)
		assert.Equal(t, c.out, res, c.in)
	}
}

func TestMonth(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{
			"a",
			false,
		},
		{
			"",
			false,
		},
		{
			"-1",
			false,
		},
		{
			"0-34",
			false,
		},
		{
			"20-15",
			false,
		},
		{
			"10,20,31,15",
			false,
		},
		{
			"feb-jan",
			false,
		},
		{
			"feb-",
			false,
		},
		{
			"-",
			false,
		},
		{
			"-jan",
			false,
		},
		{
			"dec-aug",
			false,
		},
		{
			"*",
			true,
		},
		{
			"1",
			true,
		},
		{
			"1,5,9",
			true,
		},
		{
			"1-12",
			true,
		},
		{
			"jan-dec",
			true,
		},
	}
	for _, c := range cases {
		res := isValidMonth(c.in)
		assert.Equal(t, c.out, res, c.in)
	}
}

func TestWeek(t *testing.T) {
	cases := []struct {
		in  string
		out bool
	}{
		{
			"a",
			false,
		},
		{
			"",
			false,
		},
		{
			"-1",
			false,
		},
		{
			"0-7",
			false,
		},
		{
			"20-15",
			false,
		},
		{
			"10,20,31,15",
			false,
		},
		{
			"feb-jan",
			false,
		},
		{
			"tue-",
			false,
		},
		{
			"-",
			false,
		},
		{
			"-mon",
			false,
		},
		{
			"thu-wed",
			false,
		},
		{
			"*",
			true,
		},
		{
			"1",
			true,
		},
		{
			"0,5,6",
			true,
		},
		{
			"wed,fri,sun",
			true,
		},
		{
			"0-6",
			true,
		},
		{
			"fri-sat",
			true,
		},
	}
	for _, c := range cases {
		res := isValidWeekdays(c.in)
		assert.Equal(t, c.out, res, c.in)
	}
}
