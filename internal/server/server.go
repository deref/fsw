package server

import (
	"context"

	"github.com/deref/fsw/internal/api"
	"github.com/google/uuid"
)

type Service struct {
	messages chan message
	watchers []watcher
}

type message struct {
	Thunk func() error
	Err   chan error
}

type watcher struct {
	ID   string
	Path string
}

const MaxQueueSize = 20 // TODO: Tune me.

func (svc *Service) Run(ctx context.Context) {
	svc.messages = make(chan message, MaxQueueSize)
	for {
		select {
		case <-ctx.Done():
			close(svc.messages)
			for msg := range svc.messages {
				msg.Err <- context.Canceled
			}
			return
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
	id = randomUUID()
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

func randomUUID() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uid.String()
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
