package main

import (
	"net/http"

	"github.com/deref/fsw/internal"
)

func main() {
	http.ListenAndServe(":3000", &internal.Handler{})
}
