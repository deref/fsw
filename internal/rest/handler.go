package rest

import (
	"net/http"
	"sync"

	"github.com/deref/fsw/internal/api"
)

type Handler struct {
	Service api.Service

	mx            sync.RWMutex
	subscriptions map[string]*subscription
}

type subscription struct {
	Req *http.Request
	W   http.ResponseWriter
}
