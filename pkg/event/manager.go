package event

import (
	"context"
	"golang.org/x/exp/slices"
	"mgufrone.dev/job-alerts/pkg/payload"
)

type simpleManager struct {
	handlers map[string][]int
	runners  []Handler
}

func (w *simpleManager) Publish(ctx context.Context, eventType string, py payload.Payload) error {
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
		if err := w.runners[k].Run(ctx, py); err != nil {
			return err
		}
	}
	return nil
}

func (w *simpleManager) Subscribe(eventType string, runner Handler) int {
	idx := len(w.runners)
	w.runners = append(w.runners, runner)
	if len(w.handlers[eventType]) == 0 {
		w.handlers[eventType] = []int{idx}
	} else {
		w.handlers[eventType] = append(w.handlers[eventType], idx)
	}
	return idx
}

func (w *simpleManager) Unsubscribe(eventID int) error {
	if len(w.runners) < eventID {
		return ErrEventOutOfBond
	}
	w.runners[eventID] = nil
	for k, v := range w.handlers {
		if slices.Contains(v, eventID) {
			idx := slices.Index(v, eventID)
			w.handlers[k] = slices.Delete(v, idx, 1)
		}
	}
	return nil
}

func NewSimpleManager() *simpleManager {
	return &simpleManager{handlers: map[string][]int{}}
}

//func (w *simpleManager) composeFn(handler Handler) func(ctx context.Context, eventType string, py []byte) error {
//	return func(ctx context.Context, eventType string, py []byte) error {
//		p := payload.New(py)
//		if !handler.When(p) {
//			return nil
//		}
//		spanner := sentry.StartSpan(ctx, "event_subscriber", sentry.TransctionSource(sentry.SourceComponent))
//		spanner.Description = fmt.Sprintf("event_subscriber:%s:%s", eventType, handler.Name())
//		defer spanner.Finish()
//		err := handler.Run(spanner.Context(), p)
//		if err != nil {
//			w.logger(ctx).Errorln(err)
//			spanner.Status = sentry.SpanStatusInternalError
//		}
//		return err
//	}
//}

//func (w *simpleManager) Publish(ctx context.Context, eventType string, py []byte) error {
//	span := sentry.TransactionFromContext(ctx)
//	if span == nil {
//		span = sentry.StartSpan(ctx, "event_publisher", sentry.TransctionSource(sentry.SourceComponent))
//	}
//	w.logger(ctx).Infoln("publishing event", eventType)
//	span.Description = fmt.Sprintf("event:%s", eventType)
//
//	defer func() {
//		span.Status = sentry.SpanStatusOK
//		span.Finish()
//	}()
//	handlers, ok := w.handlers[eventType]
//	if !ok || len(handlers) == 0 {
//		w.logger(ctx).Infoln("no handlers for the event. skipping", eventType)
//		return nil
//	}
//	w.logger(ctx).Infof("%d handlers will process the event", len(handlers))
//	var wg sync.WaitGroup
//	wg.Add(len(handlers))
//	for _, handler := range handlers {
//		go func(waiter *sync.WaitGroup, runner Handler, py []byte) {
//			defer waiter.Done()
//			fn := w.composeFn(runner)
//			if err := fn(span.Context(), eventType, py); err != nil {
//				span.Sampled = sentry.SampledTrue
//				w.logger(ctx).
//					WithField("payload", string(py)).
//					WithField("error", err).
//					Errorf("error at event %s for handler %s: %v", eventType, runner.Name(), err)
//				w.hub.RecoverWithContext(ctx, err)
//			}
//		}(&wg, handler, py)
//	}
//	wg.Wait()
//	return nil
//}
//
//func (w *simpleManager) Subscribe(eventType string, runner Handler) *simpleManager {
//	_, ok := w.handlers[eventType]
//	if !ok {
//		w.handlers[eventType] = []Handler{runner}
//	} else {
//		w.handlers[eventType] = append(w.handlers[eventType], runner)
//	}
//	return w
//}
