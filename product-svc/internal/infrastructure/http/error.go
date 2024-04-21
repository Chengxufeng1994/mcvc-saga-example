package common

import "errors"

var (
	// ErrInvalidParam is invalid parameter error
	ErrInvalidParam = errors.New("invalid parameter")
	// ErrUnauthorized is unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrServer is server error
	ErrServer = errors.New("server error")
)
