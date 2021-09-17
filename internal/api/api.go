package api

import "context"

type Service interface {
	CreateWatcher(context.Context, *CreateWatcherInput) (id string, err error)
	DescribeWatchers(context.Context, *DescribeWatchersInput) (*DescribeWatchersOutput, error)
	DeleteWatchers(context.Context, *DeleteWatchersInput) error
	GetEvents(context.Context, *GetEventsInput) (*GetEventsOutput, error)
	// TODO: TailEvents.
}

type CreateWatcherInput struct {
	Path string `json:"path"`
}

type CreateWatcherOutput struct {
	ID   string            `json:"id"`
	Tags map[string]string `json:"tags"`
}

type DescribeWatchersInput struct {
	IDs  []string          `json:"ids"`
	Tags map[string]string `json:"tags"`
}

type DescribeWatchersOutput struct {
	Watchers []WatcherDescription `json:"watchers"`
}

type WatcherDescription struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type DeleteWatchersInput struct {
	IDs  []string          `json:"ids"`
	Tags map[string]string `json:"tags"`
}

type GetEventsInput struct {
	WatcherID string `json:"watcherId"`
	After     string `json:"after"`
}

type GetEventsOutput struct {
	Events []Event `json:"events"`
}

type Event struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Path   string `json:"path"`
}
