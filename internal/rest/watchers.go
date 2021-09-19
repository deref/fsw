package rest

import (
	"net/http"

	"github.com/deref/fsw/internal/api"
)

type WatcherCollection struct {
	Service   api.Service
	Publisher *Publisher
}

func (coll *WatcherCollection) HandleGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	output, err := coll.Service.DescribeWatchers(ctx, &api.DescribeWatchersInput{
		// TODO: ids, tags.
	})
	writeJSONResponse(w, req, output, err)
}

func (coll *WatcherCollection) HandlePost(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var input api.CreateWatcherInput
	if err := readJSON(&input, req); err != nil {
		writeErrf(w, req, "reading json: %w", err)
		return
	}
	watcherID, err := coll.Service.CreateWatcher(ctx, &input)
	if err != nil {
		writeErr(w, req, err)
		return
	}
	w.Header().Set("Location", watcherID)
	w.WriteHeader(http.StatusCreated)
}

func (coll *WatcherCollection) Subresource(key string) interface{} {
	return &Watcher{
		Service:   coll.Service,
		Publisher: coll.Publisher,
		ID:        key,
	}
}

type Watcher struct {
	Service   api.Service
	Publisher *Publisher
	ID        string
}

func (watcher *Watcher) HandleGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	output, err := watcher.Service.DescribeWatchers(ctx, &api.DescribeWatchersInput{
		IDs: []string{watcher.ID},
	})
	if err != nil {
		writeErr(w, req, err)
		return
	}
	if len(output.Watchers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeJSONResponse(w, req, output.Watchers[0], nil)
}

func (watcher *Watcher) HandleDelete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := watcher.Service.DeleteWatchers(ctx, &api.DeleteWatchersInput{
		IDs: []string{watcher.ID},
	})
	writeErr(w, req, err)
}

func (watcher *Watcher) Subresource(key string) interface{} {
	return &EventCollection{
		Service:   watcher.Service,
		Publisher: watcher.Publisher,
		WatcherID: watcher.ID,
	}
}

type EventCollection struct {
	Service   api.Service
	Publisher *Publisher
	WatcherID string
}

func (coll *EventCollection) HandleGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	sub := coll.Publisher.subscribe(w)
	err := coll.Service.TailEvents(ctx, &api.TailEventsInput{
		WatcherID:      coll.WatcherID,
		After:          req.URL.Query().Get("after"),
		SubscriptionID: sub.id,
	})
	if err != nil {
		sub.stream.Cancel(toHttpErr(err))
	}
	sub.stream.Wait()
}
