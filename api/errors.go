package api

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidId() Error {
	return NewError(http.StatusBadRequest, "Invalid ID")
}

func ErrUnAuthenticated() Error {
	return NewError(http.StatusUnauthorized, "Forbidden")
}

func ErrUnAuthorized() Error {
	return NewError(http.StatusForbidden, "Forbidden")
}

func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "Invalid JSON Request")
}

func ErrResourceNotFound() Error {
	return NewError(http.StatusBadRequest, "Resource not found")
}
