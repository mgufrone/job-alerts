package worker

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) ToHTTPClient() *http.Client {
	return m.Called().Get(0).(*http.Client)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type MockRoundTripper struct {
	mock.Mock
}

func (r *MockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	args := r.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}
