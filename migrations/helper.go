package migrations

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrator(source *sql.DB) (gorm.Migrator, error) {
	pg, err := gorm.Open(postgres.New(postgres.Config{
		Conn: source,
	}))
	if err != nil {
		return nil, err
	}
	return pg.Migrator(), nil
}
