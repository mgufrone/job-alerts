package event

import (
	"context"
	"github.com/sirupsen/logrus"
	"mgufrone.dev/job-alerts/pkg/payload"
)

type serial struct {
	logger logrus.FieldLogger
	simpleManager
}

func newSerial(logger logrus.FieldLogger) *serial {
	mgr := &serial{logger: logger}
	mgr.handlers = map[string][]int{}
	return mgr
}

func (w *serial) Publish(ctx context.Context, eventType string, py payload.Payload) error {
	defer func() {
		if err := recover(); err != nil {
			w.logger.Errorln(err)
		}
	}()
	w.logger.Debugln("publishing event", eventType)
	handlers, ok := w.handlers[eventType]
	w.logger.Debugf("event %s has %d handlers", eventType, len(handlers))
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
		return handler.Run(ctx, py)
	}
	return nil
}
