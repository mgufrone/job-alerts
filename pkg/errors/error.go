package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound = status.Error(codes.NotFound, "proposal not found")
)

type ValidationFieldError struct {
	Field  string
	Parent error
}

func (v ValidationFieldError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Parent.Error())
}

func FieldError(name string, err error) ValidationFieldError {
	return ValidationFieldError{name, err}
}
