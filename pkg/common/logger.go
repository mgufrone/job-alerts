package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

func WithAppLogger(appName, env string) func() *logrus.Entry {
	return func() *logrus.Entry {
		return NewLogger(appName, env)
	}
}

type stderrHook struct {
	logger *logrus.Logger
}

func (h *stderrHook) Fire(entry *logrus.Entry) error {
	entry.Logger = h.logger
	return nil
}

func (h *stderrHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}
func NewLogger(appName, env string) *logrus.Entry {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{})
	logLevel := logrus.DebugLevel
	if env == "production" {
		logLevel = logrus.WarnLevel
	}
	logger.SetLevel(logLevel)
	errLogger := logrus.StandardLogger()
	errLogger.SetOutput(os.Stdout)
	logger.AddHook(&stderrHook{logger: errLogger})
	return logger.WithFields(logrus.Fields{
		"app_name": appName,
		"env":      env,
		"version":  os.Getenv("APP_VERSION"),
	})
}
