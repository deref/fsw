package main

import (
	"github.com/deref/fsw/internal"
	"github.com/deref/fsw/internal/server"
	"github.com/deref/pier"
)

func main() {
	handler := &internal.Handler{
		Service: &server.Service{},
	}
	pier.Main(handler)
}
