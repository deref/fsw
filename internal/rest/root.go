package rest

import "github.com/deref/fsw/internal/api"

type Root struct {
	Service   api.Service
	Publisher *Publisher
}

func (root *Root) Subresource(key string) interface{} {
	switch key {
	case "watchers":
		return &WatcherCollection{
			Service:   root.Service,
			Publisher: root.Publisher,
		}
	case "tags":
		return &TagCollection{
			Service: root.Service,
		}
	default:
		return nil
	}
}
