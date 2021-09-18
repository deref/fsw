package main

import (
	"context"
	"net/http"

	"github.com/deref/fsw/internal"
	"github.com/deref/fsw/internal/server"
)

func main() {
	ctx := context.Background()

	svc := &server.Service{}
	go svc.Run(ctx)

	http.ListenAndServe(":3000", &internal.Handler{
		Service: svc,
	})
}
