package httperrors

import (
	"errors"
	"net/http"
)

// Error has a status field for HTTP error codes
type Error interface {
	error
	Status() int
}

// HTTPError is an implementation of error
type HTTPError struct {
	Code int
	Err  error
}

// Error returns the error string
func (e HTTPError) Error() string {
	return e.Err.Error()
}

// Status returns the http status code
func (e HTTPError) Status() int {
	return e.Code
}

// ErrMovieNotFound wraps a movie not found case with a 404
var ErrMovieNotFound = HTTPError{http.StatusNotFound, errors.New("Movie not found")}
