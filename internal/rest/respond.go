package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deref/fsw/internal/api"
	"github.com/deref/fsw/internal/httputil"
)

func toHttpErr(err error) *httputil.Error {
	httpErr := &httputil.Error{
		Wrapped: err,
	}
	switch {
	case err == nil:
		httpErr = nil
	case errors.Is(err, api.TooBusy):
		httpErr.Status = http.StatusServiceUnavailable
	case errors.Is(err, api.NotFound):
		httpErr.Status = http.StatusNotFound
	default:
		httpErr = httputil.InternalServerError
	}
	return httpErr
}

func writeErr(w http.ResponseWriter, req *http.Request, err error) {
	httpErr := toHttpErr(err)
	if httpErr == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if httpErr.Status == http.StatusInternalServerError {
			// XXX log err
		}
		httputil.WriteErr(w, httpErr)
	}
}

func writeErrf(w http.ResponseWriter, req *http.Request, format string, v ...interface{}) {
	writeErr(w, req, fmt.Errorf(format, v...))
}

func writeJSONResponse(w http.ResponseWriter, req *http.Request, output interface{}, err error) {
	if err != nil {
		writeErr(w, req, err)
		return
	}
	writeJSON(w, req, output)
}

func readJSON(v interface{}, req *http.Request) error {
	// TODO: check content-type.
	if req.Body == nil {
		return errors.New("expected body")
	}
	dec := json.NewDecoder(req.Body)
	return dec.Decode(&v)
}

func writeJSON(w http.ResponseWriter, req *http.Request, v interface{}) {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		logf(req, "encoding json: %v", err)
	}
}

func logf(req *http.Request, format string, v ...interface{}) {
	// TODO: Include request context in log message, use contextual logger, etc.
	log.Printf(format, v...)
}
