// Command fswd stands for File System Watcher Daemon.

package main

import (
	"context"
	"net/http"

	"github.com/deref/fsw/internal/rest"
	"github.com/deref/fsw/internal/server"
)

func main() {
	ctx := context.Background()

	publisher := rest.NewPublisher()
	defer publisher.Shutdown()

	svc := &server.Service{
		Publisher: publisher,
	}

	root := &rest.Root{
		Service:   svc,
		Publisher: publisher,
	}
	handler := &rest.ResourceHandler{
		Resource: root,
	}

	go svc.Run(ctx)

	http.ListenAndServe(":3000", handler)
}
