// Code generated from Pkl module `Config`. DO NOT EDIT.
package loglevel

import (
	"encoding"
	"fmt"
)

// The level of logging for the application.
//
// - "error": Log only error level messages
// - "warn": Log error and warning messages
// - "info": Log all messages
// - "debug": Log all messages and debug information
type LogLevel string

const (
	Error LogLevel = "error"
	Warn  LogLevel = "warn"
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
)

// String returns the string representation of LogLevel
func (rcv LogLevel) String() string {
	return string(rcv)
}

var _ encoding.BinaryUnmarshaler = new(LogLevel)

// UnmarshalBinary implements encoding.BinaryUnmarshaler for LogLevel.
func (rcv *LogLevel) UnmarshalBinary(data []byte) error {
	switch str := string(data); str {
	case "error":
		*rcv = Error
	case "warn":
		*rcv = Warn
	case "info":
		*rcv = Info
	case "debug":
		*rcv = Debug
	default:
		return fmt.Errorf(`illegal: "%s" is not a valid LogLevel`, str)
	}
	return nil
}
