package apperrors

import "errors"

var (
	ErrConnectionFailed = errors.New("connection failed")
	ErrInvalidData      = errors.New("invalid data")
)
