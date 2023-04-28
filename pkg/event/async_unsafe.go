package event

import (
	"context"
	"golang.org/x/sync/errgroup"
	"mgufrone.dev/job-alerts/pkg/payload"
)

type asyncUnsafe struct {
	simpleManager
}

func newAsyncUnsafe() *asyncUnsafe {
	mgr := &asyncUnsafe{}
	mgr.handlers = map[string][]int{}
	return mgr
}

func (w *asyncUnsafe) Publish(ctx context.Context, eventType string, py payload.Payload) error {
	defer func() {
		if err := recover(); err != nil {
			// send error report
		}
	}()
	var errGroup errgroup.Group
	handlers, ok := w.handlers[eventType]
	if !ok {
		return ErrEventNotFound
	}
	for _, k := range handlers {
		if w.runners[k] == nil {
			continue
		}
		handler := w.runners[k]
		if !handler.When(py) {
			continue
		}
		errGroup.Go(func(hndl Handler) func() error {
			return func() error {
				return hndl.Run(ctx, py)
			}
		}(handler))
	}
	return errGroup.Wait()
}
