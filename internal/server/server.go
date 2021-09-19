package server

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/deref/fsw/internal/api"
	"github.com/deref/fsw/internal/gensym"
	"github.com/deref/fsw/internal/pathutil"
	"github.com/fsnotify/fsevents"
)

type Service struct {
	Publisher Publisher

	messages chan message
	events   chan []fsevents.Event
	watchers []watcher
}

type Publisher interface {
	Publish(subscriptionID string, event api.Event)
}

type message struct {
	Thunk func() error
	Err   chan error
}

type watcher struct {
	ID              string
	Path            string
	SubscriptionIDs []string
	Tags            map[string]string
}

const MaxQueueSize = 20 // TODO: Tune me.

func (svc *Service) Run(ctx context.Context) {
	svc.events = make(chan []fsevents.Event)
	svc.messages = make(chan message, MaxQueueSize)

	// TODO: Do not watch entire filesystem!
	stream := fsevents.EventStream{
		Events:  svc.events,
		Paths:   []string{"/"},
		Flags:   fsevents.FileEvents,
		Latency: time.Millisecond * 30,
	}
	stream.Start()
	defer stream.Stop()

	for {
		select {
		case <-ctx.Done():
			close(svc.messages)
			for msg := range svc.messages {
				msg.Err <- context.Canceled
			}
			return
		case events := <-svc.events:
			idBuf := make([]byte, 8)
			for _, event := range events {
				var apiEvent api.Event
				apiEvent.Path = event.Path

				binary.BigEndian.PutUint64(idBuf, event.ID)
				apiEvent.ID = hex.EncodeToString(idBuf)

				switch {
				case (event.Flags & fsevents.ItemCreated) != 0:
					apiEvent.Action = "create"
				case (event.Flags & fsevents.ItemModified) != 0:
					apiEvent.Action = "modify"
				case (event.Flags & fsevents.ItemRemoved) != 0:
					apiEvent.Action = "remove"
				default:
					// TODO: What other event types to handle?
					continue
				}

				for _, watcher := range svc.watchers {
					if pathutil.HasFilePathPrefix(event.Path, watcher.Path) {
						for _, subscriptionID := range watcher.SubscriptionIDs {
							svc.Publisher.Publish(subscriptionID, apiEvent)
						}
					}
				}
			}
		case msg := <-svc.messages:
			msg.Err <- msg.Thunk()
		}
	}
}

func (svc *Service) do(thunk func() error) error {
	err := make(chan error, 1)
	select {
	case svc.messages <- message{
		Thunk: thunk,
		Err:   err,
	}:
	default:
		err <- api.TooBusy
	}
	return <-err
}

func (svc *Service) CreateWatcher(ctx context.Context, input *api.CreateWatcherInput) (id string, err error) {
	id = gensym.RandomBase32()
	// TODO: Validate path.
	err = svc.do(func() error {
		tags := make(map[string]string, len(input.Tags))
		for name, value := range input.Tags {
			tags[name] = value
		}
		svc.watchers = append(svc.watchers, watcher{
			ID:   id,
			Path: input.Path,
			Tags: tags,
		})
		return nil
	})
	return
}

type tag struct {
	Name  string
	Value string
}

func (svc *Service) DescribeWatchers(ctx context.Context, input *api.DescribeWatchersInput) (output *api.DescribeWatchersOutput, err error) {
	var ids map[string]bool
	if input.IDs != nil {
		ids = make(map[string]bool)
		for _, id := range input.IDs {
			ids[id] = true
		}
	}

	var tags map[tag]bool
	if input.Tags != nil {
		tags = make(map[tag]bool)
		for name, value := range input.Tags {
			tags[tag{name, value}] = true
		}
	}

	err = svc.do(func() error {
		watchers := make([]api.WatcherDescription, 0, len(svc.watchers))
		for _, watcher := range svc.watchers {
			if !(ids == nil || ids[watcher.ID]) {
				continue
			}
			if tags != nil {
				match := false
				for name, value := range watcher.Tags {
					if tags[tag{name, value}] {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			desc := api.WatcherDescription{
				ID:   watcher.ID,
				Path: watcher.Path,
				Tags: make(map[string]string, len(watcher.Tags)),
			}
			for name, value := range watcher.Tags {
				desc.Tags[name] = value
			}
			watchers = append(watchers, desc)
		}
		output = &api.DescribeWatchersOutput{
			Watchers: watchers,
		}
		return nil
	})
	return
}

func (svc *Service) DeleteWatchers(ctx context.Context, input *api.DeleteWatchersInput) error {
	ids := make(map[string]bool)
	for _, id := range input.IDs {
		ids[id] = true
	}

	tags := make(map[tag]bool)
	for name, value := range input.Tags {
		tags[tag{name, value}] = true
	}

	return svc.do(func() error {
		dst := 0
		for i := 0; i < len(svc.watchers); i++ {
			keep := true
			watcher := svc.watchers[i]
			if ids[watcher.ID] {
				keep = false
			} else {
				for name, value := range watcher.Tags {
					if tags[tag{name, value}] {
						keep = false
						break
					}
				}
			}
			svc.watchers[dst] = watcher
			if keep {
				dst++
			}
		}
		svc.watchers = svc.watchers[:dst]
		return nil
	})
}

func (svc *Service) GetEvents(context.Context, *api.GetEventsInput) (*api.GetEventsOutput, error) {
	panic("TODO: GetEvents")
}

func (svc *Service) TailEvents(ctx context.Context, input *api.TailEventsInput) error {
	return svc.do(func() error {
		for i, watcher := range svc.watchers {
			if watcher.ID != input.WatcherID {
				continue
			}
			watcher.SubscriptionIDs = append(watcher.SubscriptionIDs, input.SubscriptionID)
			svc.watchers[i] = watcher
			return nil
		}
		return api.NotFound
	})
}
