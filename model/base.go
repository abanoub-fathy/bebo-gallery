package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Base is the base type like gorm.Model but it has id of type uuid instead of uint
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeSave is used to assign new UUID to the id column
// of the Base type
func (b *Base) BeforeCreate(db *gorm.DB) (err error) {
	uuid := uuid.NewV4()
	db.Statement.SetColumn("ID", uuid)
	return
}
