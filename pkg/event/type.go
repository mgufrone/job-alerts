package event

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/pkg/payload"
)

type Handler interface {
	Name() string
	Run(ctx context.Context, py payload.Payload) error
	When(py payload.Payload) bool
	SubscribeAt() []string
}

type IManager interface {
	Publish(ctx context.Context, eventType string, py payload.Payload) error
	Subscribe(eventType string, runner Handler) int
	Unsubscribe(eventID int) error
}

func AsHandler(f any) any {
	return fx.Annotate(f, fx.As(new(Handler)), fx.ResultTags(`group:"event_subscribers"`))
}

var (
	ErrEventNotFound  = fmt.Errorf("event not found or not subscribed")
	ErrEventOutOfBond = fmt.Errorf("event id is out of bond")
)
