package main

import (
	context2 "context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/internal/repositories/channel"
	"mgufrone.dev/job-alerts/internal/repositories/job"
	notification2 "mgufrone.dev/job-alerts/internal/repositories/notification"
	"mgufrone.dev/job-alerts/internal/repositories/user"
	"mgufrone.dev/job-alerts/internal/repositories/user_channel"
	"mgufrone.dev/job-alerts/internal/services/publisher"
	"mgufrone.dev/job-alerts/internal/services/upwork"
	"mgufrone.dev/job-alerts/internal/services/weworkremotely"
	"mgufrone.dev/job-alerts/internal/usecases/notification"
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/worker"
)

func main() {
	godotenv.Load()
	app := fx.New(
		fx.Provide(
			func() *logrus.Logger {
				lg := logrus.New()
				lg.SetLevel(logrus.DebugLevel)
				return lg
			},
			func(lg *logrus.Entry) logrus.FieldLogger {
				return lg
			},
			logrus.NewEntry,
			func(lg *logrus.Entry) db.Resolver {
				return func() (*gorm.DB, error) {
					db1, err := db.Open(lg)
					return db1, err
				}
			},
			db.New,
			worker.AsJobWorker(upwork.NewHandler),
			worker.AsJobWorker(weworkremotely.NewHandler),
			fx.Annotate(
				publisher.New,
				fx.ParamTags(`group:"workers"`),
			),
		),
		job.RepoModule,
		channel.RepoModule,
		user.RepoModule,
		user_channel.RepoModule,
		notification2.RepoModule,
		notification.Module,
		fx.Invoke(
			func(db *db.Instance) error {
				_ = db.DB().AutoMigrate(
					&models.Tag{},
					&models.Job{},
					&models.JobTag{},
					&models.User{},
					&models.UserChannel{},
					&models.Notification{},
					&models.Channel{},
				)
				return nil
			}, fx.Annotate(
				//	func(
				//		repo job2.QueryResolver,
				//		cmd job2.CommandResolver,
				//		usrCmd user2.CommandResolver,
				//		usrChCmd user_channel2.CommandResolver,
				//		chCmd channel2.CommandResolver,
				//		workers []worker.IWorker,
				//		db *db.Instance,
				//		logger logrus.FieldLogger,
				//	) {
				//		logger.Info("pumping data")
				//		cb := job.TagCriteria()
				//		cb = cb.
				//			Or(
				//				cb.Where(criteria.NewCondition("name", criteria.In, []string{"php", "devops"})),
				//			)
				//		_, _, _ = repo().Apply(cb).GetAndCount(context2.Background())
				//		ctx := context2.Background()
				//		command := cmd()
				//		// clean up all data
				//		db.DB().Exec("TRUNCATE jobs RESTART IDENTITY")
				//		db.DB().Exec("TRUNCATE users RESTART IDENTITY")
				//		db.DB().Exec("TRUNCATE job_tags RESTART IDENTITY")
				//		db.DB().Exec("TRUNCATE notifications RESTART IDENTITY")
				//		db.DB().Exec("TRUNCATE channels RESTART IDENTITY")
				//		db.DB().Exec("TRUNCATE user_channels RESTART IDENTITY")
				//		//db.DB().Where("id != ''").Delete(&models.Job{Job})
				//		for _, wrk := range workers {
				//			jobs, _ := wrk.Fetch(context2.Background())
				//			for _, j := range jobs {
				//				err := command.Create(ctx, j)
				//				fmt.Println("job", j.ID(), err)
				//			}
				//		}
				//		total, _ := repo().Count(ctx)
				//		var (
				//			usr user2.Entity
				//		)
				//		usr.SetAuthID("random-string")
				//		usr.SetID(uuid.New())
				//		usr.SetStatus(user2.Active)
				//		usrCmd().Create(ctx, &usr)
				//		ch := test_data.ValidChannel(true)
				//		usrCh := test_data.ValidUserChannel(&usr, "email")
				//		ch.SetUser(&usr)
				//		chCmd().Create(ctx, ch)
				//		usrChCmd().Create(ctx, usrCh)
				//		fmt.Println("inserted", total)
				//	},
				//	fx.ParamTags(
				//		`name:"source"`,
				//		`name:"source"`,
				//		`name:"source"`,
				//		`name:"source"`,
				//		`name:"source"`,
				//		`group:"workers"`,
				//	),
				//),
				//fx.Annotate(
				func(uc *notification.UseCase) error {
					return uc.Loop(context2.Background())
				},
			),
		),
	)
	ctx := context2.Background()
	app.Start(ctx)
	app.Stop(ctx)
}
