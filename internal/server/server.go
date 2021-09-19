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
}

const MaxQueueSize = 20 // TODO: Tune me.

func (svc *Service) Run(ctx context.Context) {
	svc.events = make(chan []fsevents.Event)
	svc.messages = make(chan message, MaxQueueSize)

	// TODO:
	stream := fsevents.EventStream{
		Events:  svc.events,
		Paths:   []string{"/"},
		Flags:   fsevents.FileEvents, // TODO: Include fsevents.NoDefer?
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
		svc.watchers = append(svc.watchers, watcher{
			ID:   id,
			Path: input.Path,
		})
		return nil
	})
	return
}

func (svc *Service) DescribeWatchers(ctx context.Context, input *api.DescribeWatchersInput) (output *api.DescribeWatchersOutput, err error) {
	err = svc.do(func() error {
		// TODO: handle input.IDs, etc.
		watchers := make([]api.WatcherDescription, len(svc.watchers))
		for i, watcher := range svc.watchers {
			watchers[i] = api.WatcherDescription{
				ID:   watcher.ID,
				Path: watcher.Path,
			}
		}
		output = &api.DescribeWatchersOutput{
			Watchers: watchers,
		}
		return nil
	})
	return
}

func (svc *Service) DeleteWatchers(context.Context, *api.DeleteWatchersInput) error {
	panic("TODO: DeleteWatchers")
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
