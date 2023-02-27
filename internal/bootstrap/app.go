package bootstrap

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"mgufrone.dev/job-alerts/pkg/env"
	"mgufrone.dev/job-alerts/pkg/logger"
	"os"
)

var (
	AppModule = fx.Options(
		fx.Invoke(
			func() {
				godotenv.Load()
			},
		),
		fx.Provide(
			fx.Annotate(
				logger.WithAppLogger("platform-backend", os.Getenv("APP_ENV"), env.GetOr("APP_LOG_LEVEL", "debug")),
				fx.As(new(logrus.FieldLogger)),
			),
		),
	)
)
