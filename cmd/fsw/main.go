package main

import (
	"net/http"

	"github.com/deref/fsw/internal"
	"github.com/deref/fsw/internal/server"
)

func main() {
	http.ListenAndServe(":3000", &internal.Handler{
		Service: &server.Service{},
	})
}
