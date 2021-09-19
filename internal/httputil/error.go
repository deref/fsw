package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Status  int
	Wrapped error
}

func (err *Error) Error() string {
	return fmt.Sprintf("http status %d: %v", err.Status, err.Wrapped)
}

func (err *Error) Unwrap() error {
	return err.Wrapped
}

func WriteErr(w http.ResponseWriter, err *Error) {
	w.WriteHeader(err.Status)
	w.Write([]byte(fmt.Sprintf("status %d: %v\n", err.Status, err.Wrapped)))
}

var InternalServerError = &Error{
	Status:  http.StatusInternalServerError,
	Wrapped: errors.New("internal server error"),
}
