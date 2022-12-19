package sentry

import (
	"github.com/getsentry/sentry-go"
	"os"
)

func Initialize() error {
	return sentry.Init(sentry.ClientOptions{Dsn: os.Getenv("SENTRY_DSN")})
}
func AttachAppTags(appName string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("app_version", os.Getenv("APP_VERSION"))
		scope.SetTag("app_env", os.Getenv("APP_ENV"))
		scope.SetTag("app_name", appName)
	})
}
