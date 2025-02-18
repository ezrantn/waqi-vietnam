package main

import (
	"fmt"
)

type ErrorType string

const (
	ErrInvalidCity     ErrorType = "INVALID_CITY"
	ErrAPIUnavailable  ErrorType = "API_UNAVAILABLE"
	ErrInvalidResponse ErrorType = "INVALID_RESPONSE"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
