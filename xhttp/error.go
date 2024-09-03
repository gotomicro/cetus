package xhttp

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrBadRequest          = errors.New(http.StatusText(http.StatusBadRequest))
	ErrUnauthorized        = errors.New(http.StatusText(http.StatusUnauthorized))
	ErrForbidden           = errors.New(http.StatusText(http.StatusForbidden))
	ErrNotFound            = errors.New(http.StatusText(http.StatusNotFound))
	ErrInternalServerError = errors.New(http.StatusText(http.StatusInternalServerError))
)

func ErrFromCode(code int) error {
	switch code {
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return ErrInternalServerError
	}
}
