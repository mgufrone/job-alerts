package publisher

import (
	"context"
	"errors"
	"mgufrone.dev/job-alerts/internal/domain/notification"
	"mgufrone.dev/job-alerts/internal/domain/user_channel"
)

type Publisher interface {
	Name() string
	Publish(ctx context.Context, entity *notification.Entity) error
	SetReceiver(entity *user_channel.Entity)
}

type Collection struct {
	publishers map[string]Publisher
}

func New(publishers []Publisher) *Collection {
	coll := &Collection{publishers: map[string]Publisher{}}
	for _, c := range publishers {
		if _, ok := coll.publishers[c.Name()]; ok {
			continue
		}
		coll.publishers[c.Name()] = c
	}
	return coll
}

func (c *Collection) Get(name string) (Publisher, error) {
	if p, ok := c.publishers[name]; ok {
		return p, nil
	}
	return nil, errors.New("publisher not found")
}

func (c *Collection) Clear() *Collection {
	c.publishers = map[string]Publisher{}
	return c
}
