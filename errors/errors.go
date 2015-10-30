package errors

import (
	"errors"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func (e HTTPError) Status() int {
	return e.Code
}

var ErrMovieNotFound = HTTPError{http.StatusNotFound, errors.New("Movie not found")}
