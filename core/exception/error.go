package exception

import (
	"errors"
	"fmt"
)

type HttpError struct {
	Status  int
	Message string
	Cause   error
}

func (e *HttpError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%d %s: %v", e.Status, e.Message, e.Cause)
	}
	return fmt.Sprintf("%d %s", e.Status, e.Message)
}

func New(status int, msg string) *HttpError { return &HttpError{Status: status, Message: msg} }
func Wrap(status int, msg string, err error) *HttpError {
	return &HttpError{Status: status, Message: msg, Cause: err}
}

func FromRecover(r any) error {
	switch v := r.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return fmt.Errorf("panic: %v", v)
	}
}
