package utils

import (
	"errors"
	"net/http"
)

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	if e.Message == "" {
		return "not found"
	}
	return e.Message
}

type InsertionError struct {
	Message string
}

func (e *InsertionError) Error() string {
	if e.Message == "" {
		return "insertion error"
	}
	return e.Message
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	if e.Message == "" {
		return "bad request"
	}
	return e.Message
}

type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	if e.Message == "" {
		return "full authentication is required to access this resource"
	}
	return e.Message
}

type AccessDeniedError struct {
	Message string
}

func (e *AccessDeniedError) Error() string {
	if e.Message == "" {
		return "access denied"
	}
	return e.Message
}

type ServerError struct {
	Message string
}

func (e *ServerError) Error() string {
	if e.Message == "" {
		return "server error"
	}
	return e.Message
}

func errorStatus(err error) int {
	var notFoundError *NotFoundError
	var insertionError *InsertionError
	var badRequestError *BadRequestError
	var authenticationError *AuthenticationError
	var accessDeniedError *AccessDeniedError
	var serverError *ServerError
	switch {
	case errors.As(err, &notFoundError):
		return http.StatusNotFound
	case errors.As(err, &insertionError):
		return http.StatusConflict
	case errors.As(err, &badRequestError):
		return http.StatusBadRequest
	case errors.As(err, &authenticationError):
		return http.StatusUnauthorized
	case errors.As(err, &accessDeniedError):
		return http.StatusForbidden
	case errors.As(err, &serverError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), errorStatus(err))
}
