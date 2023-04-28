package bootstrap

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/repositories/channel"
	"mgufrone.dev/job-alerts/internal/repositories/job"
	notification2 "mgufrone.dev/job-alerts/internal/repositories/notification"
	"mgufrone.dev/job-alerts/internal/repositories/user"
	"mgufrone.dev/job-alerts/internal/repositories/user_channel"
	"mgufrone.dev/job-alerts/internal/services/publisher"
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/logger"
)

var (
	DBModule = fx.Options(
		fx.Provide(
			func(lg logrus.FieldLogger) db.Resolver {
				return func() (*gorm.DB, error) {
					db1, err := db.Open(lg.(*logrus.Entry))
					return db1, err
				}
			},
			db.New,
		),
	)
	UtilityModule = fx.Options(
		fx.Provide(
			func() []publisher.Publisher {
				return nil
			},
			publisher.New,
		),
	)
	RepoModule = fx.Options(
		job.RepoModule,
		channel.RepoModule,
		user.RepoModule,
		user_channel.RepoModule,
		notification2.RepoModule,
	)
	AppModule = fx.Options(
		logger.Module,
		DBModule,
		UtilityModule,
		RepoModule,
	)
)
