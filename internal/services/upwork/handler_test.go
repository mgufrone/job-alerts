package upwork

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	worker2 "mgufrone.dev/job-alerts/pkg/worker"
	wrapper "mgufrone.dev/job-alerts/pkg/worker_wrapper"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestHandler_Fetch(t *testing.T) {
	t.Parallel()
	mockXML, _ := os.ReadFile("./rss_sample.xml")
	res := &http.Response{
		Body:       io.NopCloser(strings.NewReader(string(mockXML))),
		StatusCode: http.StatusOK,
	}
	handler := NewHandler(mockClient(res, nil))
	jbs, err := handler.Fetch(context.Background())
	assert.Len(t, jbs, 10)
	assert.Nil(t, err)
	assert.Equal(t, "hourly", jbs[0].JobType())
	assert.Equal(t, "Back-End Development", jbs[0].Role())
	//assert.Equal(t, "Credit Repair Cloud", jbs[0].CompanyName())
	assert.Equal(t, "-", jbs[0].CompanyURL())
	assert.Equal(t, "https://www.upwork.com/jobs/Senior-Backend-Engineer_%7E01da1f8d048472d18f?source=rss", jbs[0].JobURL())
	assert.Equal(t, []string{"Golang", "Microservice", "PostgreSQL", "JavaScript", "API", "MongoDB", "Software Architecture & Design", "RESTful API", "Database"}, jbs[0].Tags())
	assert.Equal(t, WorkerName, jbs[0].Source())
	for _, j := range jbs {
		min, max := j.Salary()[0], j.Salary()[1]
		if j.JobType() == "hourly" {
			assert.True(t, min < max)
		}
		if j.JobType() == "unknown" {
			assert.Equal(t, 0.01, min)
			assert.Equal(t, 0.01, max)
		}
		if j.JobType() == "fixed-price" {
			assert.True(t, min == max)
			assert.Greater(t, max, 0.00)
		}
	}
}

func mockClient(res *http.Response, err error) worker2.IHTTPClient {
	transport := &worker2.MockRoundTripper{}
	transport.On("RoundTrip", mock.Anything).
		Return(res, err)
	cli := &http.Client{}
	cli.Transport = transport
	return wrapper.NewHTTPClient(cli)
}

func TestHandler_FetchFail(t *testing.T) {
	t.Parallel()
	res := &http.Response{
		Body:       io.NopCloser(strings.NewReader(``)),
		StatusCode: http.StatusInternalServerError,
	}
	handler := NewHandler(mockClient(res, nil))
	jbs, err := handler.Fetch(context.Background())
	assert.Nil(t, jbs)
	assert.NotNil(t, err)
}
