package sse

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deref/fsw/internal/httputil"
)

type Stream struct {
	ResponseWriter http.ResponseWriter

	wroteHeader bool
	err         chan *httputil.Error
	data        chan []byte
	done        chan struct{}
}

var StreamBuffer = 100
var KeepAliveTimeout = 10 * time.Second

func (stream *Stream) Init() {
	stream.err = make(chan *httputil.Error, 1)
	stream.done = make(chan struct{})
	stream.data = make(chan []byte, StreamBuffer)
}

func (stream *Stream) Run() {
	defer close(stream.done)

	w := stream.ResponseWriter

	for {
		select {
		case err := <-stream.err:
			stream.writeHeader(err)
			return
		case <-time.After(KeepAliveTimeout):
			stream.writeHeader(nil)
			if _, err := fmt.Fprintf(w, "; %s\n\n", time.Now().UTC()); err != nil {
				return
			}
		case data := <-stream.data:
			stream.writeHeader(nil)
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return
			}
		}
	}
}

func (stream *Stream) writeHeader(err *httputil.Error) {
	if stream.wroteHeader {
		return
	}
	stream.wroteHeader = true

	w := stream.ResponseWriter
	if err != nil {
		httputil.WriteErr(w, err)
		return
	}
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
}

func (stream *Stream) Cancel(err *httputil.Error) {
	select {
	case stream.err <- err:
	default:
	}
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
		stream.Cancel(nil)
	}
}
