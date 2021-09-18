package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deref/fsw/internal/api"
)

type Handler struct {
	Service api.Service
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	switch req.URL.Path {

	case "/watchers":
		switch req.Method {

		case "GET":
			output, err := h.Service.DescribeWatchers(ctx, &api.DescribeWatchersInput{
				// TODO: ids, tags.
			})
			writeJSONResponse(w, req, output, err)

		case "POST":
			var input api.CreateWatcherInput
			if err := readJSON(&input, req); err != nil {
				writeErrf(w, req, "reading json: %w", err)
				return
			}
			output, err := h.Service.CreateWatcher(ctx, &input)
			writeJSONResponse(w, req, output, err)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func writeErr(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, api.TooBusy):
		status = http.StatusServiceUnavailable
	}
	w.WriteHeader(status)
	if status == http.StatusInternalServerError {
		logf(req, "reading json: %v", err)
	} else {
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
