package errors

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClientError struct {
	content *http.Response
	req     *http.Request
}

func NewClientError(content *http.Response, req *http.Request) *ClientError {
	return &ClientError{
		content, req,
	}
}

func (s *ClientError) Error() string {
	var body []byte
	res := map[string]interface{}{
		"status":  s.content.StatusCode,
		"request": s.req.URL.String(),
	}
	if s.content != nil && s.content.Body != nil {
		body, _ = ioutil.ReadAll(s.content.Body)
		res["body"] = string(body)
	}
	return fmt.Sprintf("clientError: %v", res)
}
