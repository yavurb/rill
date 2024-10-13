package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/yavurb/rill/internal/broadcasts/domain"
)

type EchoLogger struct {
	echoLogger echo.Logger
}

func NewEchoLogger(echoLogger echo.Logger) domain.Logger {
	return echoLogger
}
