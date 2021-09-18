package main

import (
	"context"
	"net/http"

	"github.com/deref/fsw/internal/rest"
	"github.com/deref/fsw/internal/server"
)

func main() {
	ctx := context.Background()

	svc := &server.Service{}

	root := &rest.Root{
		Service: svc,
	}
	handler := &rest.ResourceHandler{
		Resource: root,
	}
	svc.Publisher = root

	go svc.Run(ctx)

	http.ListenAndServe(":3000", handler)
}
