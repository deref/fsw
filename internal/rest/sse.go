package rest

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/deref/fsw/internal/api"
	"github.com/deref/fsw/internal/sse"
	"github.com/deref/fsw/internal/uuidutil"
)

type Publisher struct {
	mx   sync.Mutex
	subs map[string]*Subscription
}

type Subscription struct {
	id     string
	stream *sse.Stream
}

func NewPublisher() *Publisher {
	return &Publisher{
		subs: make(map[string]*Subscription),
	}
}

func (pub *Publisher) Shutdown() {
	pub.mx.Lock()
	defer pub.mx.Unlock()
	for _, sub := range pub.subs {
		sub.stream.Cancel()
	}
}

func (pub *Publisher) subscribe(w http.ResponseWriter) *Subscription {
	pub.mx.Lock()
	defer pub.mx.Unlock()

	sub := &Subscription{
		id: uuidutil.RandomString(),
		stream: &sse.Stream{
			ResponseWriter: w,
		},
	}
	pub.subs[sub.id] = sub

	sub.stream.Init()
	go func() {
		sub.stream.Run()
		pub.removeStream(sub.id)
	}()

	return sub
}

func (pub *Publisher) Publish(subscriptionID string, event api.Event) {
	pub.mx.Lock()
	defer pub.mx.Unlock()
	sub := pub.subs[subscriptionID]
	if sub == nil {
		return
	}
	data, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	sub.stream.SendData(data)
}

func (pub *Publisher) removeStream(id string) {
	pub.mx.Lock()
	defer pub.mx.Unlock()
	delete(pub.subs, id)
}
