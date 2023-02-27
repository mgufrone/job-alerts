package db

import "gorm.io/gorm"

type Resolver func() (*gorm.DB, error)
type Instance struct {
	db       *gorm.DB
	resolver func() (*gorm.DB, error)
}

func New(resolver Resolver) *Instance {
	return &Instance{resolver: resolver}
}

func (i *Instance) DB() *gorm.DB {
	if i.db == nil {
		db, err := i.resolver()
		if err != nil {
			panic(err)
		}
		i.db = db
	}
	return i.db
}
