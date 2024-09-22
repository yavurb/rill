package domain

import "errors"

var (
	ErrUnknown               = errors.New("unknown error")
	ErrBroadcastNotFound     = errors.New("broadcast not found")
	ErrBroadcastInvalidEvent = errors.New("broadcast invalid event")
)
