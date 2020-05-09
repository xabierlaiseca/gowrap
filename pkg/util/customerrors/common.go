package customerrors

import (
	"errors"
	"fmt"
)

func Error(message string) error {
	return &SimpleError{message: message}
}

func ErrorWithCause(cause error, message string) error {
	return &chainedError{SimpleError: &SimpleError{message: message}, cause: cause}
}

func Errorf(format string, a ...interface{}) error {
	return Error(fmt.Sprintf(format, a...))
}

func NotFound() error {
	return &notFoundError{}
}

func IsNotFound(err error) bool {
	return errors.Is(err, &notFoundError{})
}

type SimpleError struct {
	message string
}

func (gwe *SimpleError) Error() string {
	return gwe.message
}

type chainedError struct {
	*SimpleError
	cause error
}

func (gwe *chainedError) Unwrap() error {
	return gwe.cause
}

type notFoundError struct{}

func (*notFoundError) Error() string {
	return "resource not found"
}
