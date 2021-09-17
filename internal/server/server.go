package server

import (
	"context"
	"sync"

	"github.com/deref/fsw/internal/api"
	"github.com/google/uuid"
)

type Service struct {
	mx       sync.RWMutex
	watchers []watcher
}

type watcher struct {
	ID   string
	Path string
}

func (svc *Service) CreateWatcher(ctx context.Context, input *api.CreateWatcherInput) (id string, err error) {
	svc.mx.Lock()
	defer svc.mx.Unlock()

	id = randomUUID()
	// TODO: Validate path.
	svc.watchers = append(svc.watchers, watcher{
		ID:   id,
		Path: input.Path,
	})
	return id, nil
}

func randomUUID() string {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uid.String()
}

func (svc *Service) DescribeWatchers(ctx context.Context, input *api.DescribeWatchersInput) (*api.DescribeWatchersOutput, error) {
	svc.mx.RLock()
	defer svc.mx.RUnlock()

	// TODO: handle input.IDs, etc.
	watchers := make([]api.WatcherDescription, len(svc.watchers))
	for i, watcher := range svc.watchers {
		watchers[i] = api.WatcherDescription{
			ID:   watcher.ID,
			Path: watcher.Path,
		}
	}
	return &api.DescribeWatchersOutput{
		Watchers: watchers,
	}, nil
}

func (svc *Service) DeleteWatchers(context.Context, *api.DeleteWatchersInput) error {
	panic("TODO: DeleteWatchers")
}

func (svc *Service) GetEvents(context.Context, *api.GetEventsInput) (*api.GetEventsOutput, error) {
	panic("TODO: GetEvents")
}
