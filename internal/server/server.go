package server

import (
	"context"

	"github.com/deref/fsw/internal/api"
)

type Service struct {
	watchers []watcher
}

type watcher struct {
	ID   string
	Path string
}

func (svc *Service) CreateWatcher(ctx context.Context, input *api.CreateWatcherInput) (id string, err error) {
	panic("TODO: CreateWatcher")
}

func (svc *Service) DescribeWatchers(ctx context.Context, input *api.DescribeWatchersInput) (*api.DescribeWatchersOutput, error) {
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
