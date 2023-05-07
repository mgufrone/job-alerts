package migrations

import (
	"database/sql"
	"github.com/mgufrone/go-utils/try"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/internal/models"
)

func init() {
	goose.AddMigrationNoTx(upJobTitle, downJobTitle)
}

func upJobTitle(source *sql.DB) error {
	var (
		mgr gorm.Migrator
	)
	return try.RunOrError(func() (err error) {
		mgr, err = migrator(source)
		return
	}, func() error {
		if !mgr.HasColumn(&models.Job{}, "Title") {
			return mgr.AddColumn(&models.Job{}, "Title")
		}
		return nil
		// TODO: do your migration here
		return nil
	})
}

func downJobTitle(source *sql.DB) error {
	// remove this command to activate down migration
	//	var (
	//		mgr gorm.Migrator
	//	)
	/** remove this command to activate down migration
	return try.RunOrError(func() (err error){
		mgr, err = migrator(source)
		return
	}, func() error {
		// TODO: do your migration here
		return nil
	})
	*/
	return nil
}
