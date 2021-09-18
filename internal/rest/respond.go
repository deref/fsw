package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deref/fsw/internal/api"
)

func writeErr(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError
	switch {
	case err == nil:
		status = http.StatusOK
	case errors.Is(err, api.TooBusy):
		status = http.StatusServiceUnavailable
	}
	w.WriteHeader(status)
	if status == http.StatusInternalServerError {
		logf(req, "reading json: %v", err)
	} else if err != nil {
		w.Write([]byte(err.Error() + "\n"))
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
