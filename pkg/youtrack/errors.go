package youtrack

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("youtrack: unauthorized")
	ErrNotFound     = errors.New("youtrack: not found")
	ErrForbidden    = errors.New("youtrack: forbidden")
)

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("youtrack: HTTP %d: %s", e.StatusCode, e.Message)
}
