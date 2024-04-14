package domain

import "errors"

var (
	ErrUnknown  = errors.New("unknown error")
	ErrNotFound = errors.New("broadcast not found")
)
