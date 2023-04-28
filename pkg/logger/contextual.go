package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

var (
	loggerContext = loggerCtx{}
)

type loggerCtx struct {
}
type FnLoggerCtx func(ctx context.Context) logrus.FieldLogger
type Contextual struct {
	logger logrus.FieldLogger
}

func NewContextual(logger logrus.FieldLogger) *Contextual {
	return &Contextual{logger: logger}
}

func (t *Contextual) WithContext(ctx context.Context) (context.Context, logrus.FieldLogger) {
	logger := t.logger.(*logrus.Entry).WithContext(ctx)
	ctx = context.WithValue(ctx, loggerContext, logger)
	return ctx, logger
}

// Logger will try to guess if context is already have logger attached to the context
// otherwise, it will create copy logger from Contextual
func (t *Contextual) Logger(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value(loggerContext).(logrus.FieldLogger); ok && logger != nil {
		return logger
	}
	_, logger := t.WithContext(ctx)
	return logger
}
