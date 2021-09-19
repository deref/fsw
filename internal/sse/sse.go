package sse

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Stream struct {
	ResponseWriter http.ResponseWriter

	ctx    context.Context
	cancel func()
	data   chan []byte
	done   chan struct{}
}

var StreamBuffer = 100
var KeepAliveTimeout = 10 * time.Second

func (stream *Stream) Init() {
	ctx := context.Background()
	stream.ctx, stream.cancel = context.WithCancel(ctx)
	stream.done = make(chan struct{})

	stream.data = make(chan []byte, StreamBuffer)
}

func (stream *Stream) Run() {
	defer close(stream.done)

	w := stream.ResponseWriter
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")

	for {
		select {
		case <-stream.ctx.Done():
			return
		case <-time.After(KeepAliveTimeout):
			if _, err := fmt.Fprintf(w, "; %s\n\n", time.Now().UTC()); err != nil {
				return
			}
		case data := <-stream.data:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
		}
	}
}

func (stream *Stream) Cancel() {
	stream.cancel()
}

func (stream *Stream) Wait() {
	<-stream.done
}

// Data may not contain newlines or non-printable characters.
func (stream *Stream) SendData(data []byte) {
	select {
	case stream.data <- data:
	default:
		// Slow consumer.
		stream.cancel()
	}
}
