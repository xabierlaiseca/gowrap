package customerrors

import "fmt"

func New(message string) error {
	return &SimpleError{message: message}
}

func NewWithCause(cause error, message string) error {
	return &chainedError{SimpleError: &SimpleError{message: message}, cause: cause}
}

func Errorf(format string, a ...interface{}) error {
	return New(fmt.Sprintf(format, a...))
}

func NotFound(message string) *NotFoundError {
	return &NotFoundError{
		SimpleError: &SimpleError{message: message},
	}
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

type NotFoundError struct {
	*SimpleError
}
