package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Message string `json:"message"`
}

func (e HTTPError) InternalServerError() error {
	return echo.NewHTTPError(http.StatusInternalServerError, e.Message)
}

func (e HTTPError) BadRequest() error {
	return echo.NewHTTPError(http.StatusBadRequest, e.Message)
}

func (e HTTPError) NotFound() error {
	return echo.NewHTTPError(http.StatusNotFound, e.Message)
}

func (e HTTPError) Unauthorized() error {
	return echo.NewHTTPError(http.StatusUnauthorized, e.Message)
}

func (e HTTPError) Forbidden() error {
	return echo.NewHTTPError(http.StatusForbidden, e.Message)
}

func (e HTTPError) Conflict() error {
	return echo.NewHTTPError(http.StatusConflict, e.Message)
}

func (e HTTPError) ErrUnprocessableEntity() error {
	err := echo.ErrUnprocessableEntity

	err.Message = e.Message

	return err
}
