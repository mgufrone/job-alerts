package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`
}

func (b *Entity) BeforeCreate(scope *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

type EntityTimestamp struct {
	CreatedAt time.Time `gorm:"autoCreateTime;type:datetime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;type:datetime"`
}
