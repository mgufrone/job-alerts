package db

import (
	"fmt"
	"github.com/aklinkert/go-gorm-logrus-logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"mgufrone.dev/job-alerts/pkg/env"
	"time"
)

func logger(lg *logrus.Entry) logger2.Interface {
	if lg == nil {
		return nil
	}
	return gormlogruslogger.
		NewGormLogrusLogger(lg.WithField("component", "gorm"), 100*time.Millisecond)
}
func Open(entry *logrus.Entry) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=%s",
		env.Get("DATABASE_HOST"),
		env.Get("DATABASE_USER"),
		env.Get("DATABASE_PASS"),
		env.Get("DATABASE_DB"),
		env.Get("DATABASE_PORT"),
		env.GetOr("TIMEZONE", "UTC"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
		Logger:                                   logger(entry),
	})
}

func MustOpen() *gorm.DB {
	db, err := Open(logrus.StandardLogger().WithField("driver", "gorm"))
	if err != nil {
		panic(err)
	}
	return db
}
