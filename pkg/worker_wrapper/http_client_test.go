package wrapper

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mgufrone.dev/job-alerts/pkg/worker"
	"net/http"
	"testing"
)

func mockRoundTripper(fail bool, message string) *worker.MockRoundTripper {
	rt := &worker.MockRoundTripper{}
	var (
		res *http.Response
		err error
	)
	if fail {
		err = errors.New(message)
	} else {
		res = &http.Response{StatusCode: http.StatusOK}
	}
	rt.On("RoundTrip", mock.Anything).Return(res, err)
	return rt
}

func TestNewHttpClient(t *testing.T) {
	t.Parallel()
	testCases := []*http.Client{
		nil,
		{},
	}
	for _, c := range testCases {
		cli := NewHTTPClient(c)
		assert.Equal(t, c == nil, cli == nil)
	}
}
func TestHttpClient_Do(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		fail bool
		msg  string
	}{
		{true, "error"},
		{false, ""},
	}
	for _, tc := range testCases {
		cli := &http.Client{}
		cli.Transport = mockRoundTripper(tc.fail, tc.msg)
		agent := NewHTTPClient(cli)
		req, _ := http.NewRequest("GET", "https://somedomain.host/something", nil)
		res, err := agent.Do(req)
		if tc.fail {
			assert.Nil(t, res, tc.fail)
			assert.NotNil(t, err, tc.fail)
			continue
		}
		assert.NotNil(t, res)
		assert.Nil(t, err)
	}
}
func TestHttpClient_PredefinedAgent(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in string
	}{
		{""},
		{"abcd"},
		{RandomUserAgent()},
	}
	for _, tc := range testCases {
		cli := &http.Client{}
		cli.Transport = mockRoundTripper(true, "not ok")
		agent := NewHTTPClient(cli)
		req, _ := http.NewRequest("GET", "https://somedomain.host/something", nil)
		fmt.Println("setting header", tc.in)
		req.Header.Add("user-agent", tc.in)
		agent.Do(req)
		assert.NotEmpty(t, req.Header.Get("user-agent"))
		if tc.in != "" {
			assert.Equal(t, tc.in, req.Header.Get("user-agent"))
		}
	}
}
