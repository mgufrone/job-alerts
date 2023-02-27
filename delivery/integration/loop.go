package main

import (
	context2 "context"
	"fmt"
	gormlogruslogger "github.com/aklinkert/go-gorm-logrus-logger"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	job2 "mgufrone.dev/job-alerts/internal/domain/job"
	"mgufrone.dev/job-alerts/internal/models"
	"mgufrone.dev/job-alerts/internal/repositories/job"
	"mgufrone.dev/job-alerts/internal/services/upwork"
	"mgufrone.dev/job-alerts/internal/services/weworkremotely"
	"mgufrone.dev/job-alerts/pkg/db"
	"mgufrone.dev/job-alerts/pkg/infrastructure/criteria"
	"mgufrone.dev/job-alerts/pkg/worker"
	"time"
)

func main() {
	app := fx.New(
		fx.Provide(
			func() *logrus.Logger {
				lg := logrus.New()
				lg.SetLevel(logrus.InfoLevel)
				return lg
			},
			logrus.NewEntry,
			func(lg *logrus.Entry) db.Resolver {
				return func() (*gorm.DB, error) {
					lite := sqlite.Open("test.db")
					return gorm.Open(lite, &gorm.Config{
						QueryFields: true,
						Logger:      gormlogruslogger.NewGormLogrusLogger(lg, time.Second).LogMode(logger.Info),
					})
				}
			},
			db.New,
			worker.AsJobWorker(upwork.NewHandler),
			worker.AsJobWorker(weworkremotely.NewHandler),
		),
		job.RepoModule,
		fx.Invoke(
			func(db *db.Instance) error {
				return db.DB().AutoMigrate(&models.JobTag{}, &models.Tag{}, &models.Job{})
			}, fx.Annotate(
				func(repo job2.QueryResolver, cmd job2.CommandResolver, workers []worker.IWorker, db *db.Instance) {
					cb := job.TagCriteria()
					cb = cb.
						Or(
							cb.Where(criteria.NewCondition("name", criteria.In, []string{"php", "devops"})),
						)
					_, _, _ = repo().Apply(cb).GetAndCount(context2.Background())
					ctx := context2.Background()
					command := cmd()
					db.DB().Where("id != ''").Delete(&models.Job{})
					for _, wrk := range workers {
						jobs, _ := wrk.Fetch(context2.Background())
						for _, j := range jobs {
							err := command.Create(ctx, j)
							fmt.Println("job", j.ID(), err)
						}
					}
					total, _ := repo().Count(ctx)
					fmt.Println("inserted", total)
				},
				fx.ParamTags(`name:"source"`, `name:"source"`, `group:"workers"`),
			),
		),
	)
	ctx := context2.Background()
	app.Start(ctx)
	app.Stop(ctx)
}
