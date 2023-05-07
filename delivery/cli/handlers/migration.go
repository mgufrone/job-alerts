package handlers

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"mgufrone.dev/job-alerts/pkg/db"
	"text/template"
)

var tpl = template.Must(template.New("goose.go-migration").Parse(`package migrations

import (
	"database/sql"
	"github.com/mgufrone/go-utils/try"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationNoTx(up{{.CamelName}}, down{{.CamelName}})
}

func up{{.CamelName}}(source *sql.DB) error {
	var (
		mgr gorm.Migrator
	)
	return try.RunOrError(func() (err error){
		mgr, err = migrator(source)
		return
	}, func() error {
		// TODO: do your migration here 
		return nil
	})
}

func down{{.CamelName}}(source *sql.DB) error {
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
`))

type Migration struct {
	db *db.Instance
}
type MigrationContext struct {
	DryRun bool
}

func NewMigration(db *db.Instance) *Migration {
	return &Migration{db: db}
}

func (m *Migration) getDb(ctx MigrationContext) (*sql.DB, error) {
	ormSource := m.db.DB().Session(&gorm.Session{
		DryRun: ctx.DryRun,
	})
	orm, err := ormSource.DB()
	if err != nil {
		return nil, err
	}
	driverName := ormSource.Name()
	goose.SetDialect(driverName)
	goose.SetVerbose(ctx.DryRun)
	return orm, nil
}
func (m *Migration) Up(ctx MigrationContext) error {
	driver, err := m.getDb(ctx)
	if err != nil {
		return err
	}
	return goose.Up(driver, "migrations")
}

func (m *Migration) Down(ctx MigrationContext) error {
	driver, err := m.getDb(ctx)
	if err != nil {
		return err
	}
	return goose.Down(driver, "migrations")
}

func (m *Migration) Create(name string) error {
	driver, err := m.getDb(MigrationContext{DryRun: false})
	if err != nil {
		return err
	}
	return goose.CreateWithTemplate(driver, "migrations", tpl, name, "go")
}

func (m *Migration) Status() error {
	driver, err := m.getDb(MigrationContext{DryRun: false})
	if err != nil {
		return err
	}
	return goose.Status(driver, "migrations")
}
