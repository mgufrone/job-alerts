package migrations

import (
	"database/sql"
	"github.com/mgufrone/go-utils/try"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/models"
)

func init() {
	goose.AddMigrationNoTx(upInitialMigration, downInitialMigration)
}

func upInitialMigration(source *sql.DB) error {
	var (
		mgr gorm.Migrator
	)
	return try.RunOrError(func() (err error) {
		mgr, err = migrator(source)
		return
	}, func() error {
		if !mgr.HasTable(&models.Job{}) {
			return mgr.CreateTable(&models.Job{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.Tag{}) {
			return mgr.CreateTable(&models.Tag{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.User{}) {
			return mgr.CreateTable(&models.Job{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.UserChannel{}) {
			return mgr.CreateTable(&models.UserChannel{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.Notification{}) {
			return mgr.CreateTable(&models.Notification{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.Channel{}) {
			return mgr.CreateTable(&models.Channel{})
		}
		return nil
	}, func() error {
		if !mgr.HasTable(&models.JobTag{}) {
			return mgr.CreateTable(&models.JobTag{})
		}
		return nil
	})
}

func downInitialMigration(source *sql.DB) error {
	var (
		mgr gorm.Migrator
	)
	return try.RunOrError(func() (err error) {
		mgr, err = migrator(source)
		return
	}, func() error {
		if mgr.HasTable(&models.Job{}) {
			return mgr.DropTable(&models.Job{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.Tag{}) {
			return mgr.DropTable(&models.Tag{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.User{}) {
			return mgr.DropTable(&models.Job{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.UserChannel{}) {
			return mgr.DropTable(&models.UserChannel{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.Notification{}) {
			return mgr.DropTable(&models.Notification{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.Channel{}) {
			return mgr.DropTable(&models.Channel{})
		}
		return nil
	}, func() error {
		if mgr.HasTable(&models.JobTag{}) {
			return mgr.DropTable(&models.JobTag{})
		}
		return nil
	})
}
