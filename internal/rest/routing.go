package rest

import (
	"net/http"
	"strings"
)

type ResourceHandler struct {
	Resource interface{}
}

func (h *ResourceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resource := Route(h.Resource, req.URL.Path)
	if resource == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch req.Method {
	case "GET":
		if handler, ok := resource.(GetHandler); ok {
			handler.HandleGet(w, req)
			return
		}
	case "POST":
		if handler, ok := resource.(PostHandler); ok {
			handler.HandlePost(w, req)
			return
		}
	case "DELETE":
		if handler, ok := resource.(DeleteHandler); ok {
			handler.HandleDelete(w, req)
			return
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

type ParentResource interface {
	Subresource(key string) interface{}
}

type GetHandler interface {
	HandleGet(w http.ResponseWriter, req *http.Request)
}

type PostHandler interface {
	HandlePost(w http.ResponseWriter, req *http.Request)
}

type DeleteHandler interface {
	HandleDelete(w http.ResponseWriter, req *http.Request)
}

func Route(resource interface{}, path string) interface{} {
	// XXX Properly handle trailing slashes.
	parts := strings.Split(path, "/")[1:]
	for len(parts) > 0 {
		if parts[0] == "" {
			return resource
		}
		parent, ok := resource.(ParentResource)
		if !ok {
			return nil
		}
		resource = parent.Subresource(parts[0])
		parts = parts[1:]
	}
	return resource
}
