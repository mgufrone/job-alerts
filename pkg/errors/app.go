package errors

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (a AppError) Error() string {
	return fmt.Sprintf("%d: %s", a.Code, a.Message)
}
