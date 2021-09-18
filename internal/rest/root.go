package rest

import "github.com/deref/fsw/internal/api"

type Root struct {
	Service api.Service
}

func (root *Root) Subresource(key string) interface{} {
	switch key {
	case "watchers":
		return &WatcherCollection{
			Service: root.Service,
		}
	case "tags":
		return &TagCollection{
			Service: root.Service,
		}
	default:
		return nil
	}
}

// XXX don't put this on root
func (root *Root) Publish(subscriptionID string, event api.Event) {
}
