package model

import uuid "github.com/satori/go.uuid"

// Gallery is the container for images we will add
type Gallery struct {
	Base
	Title  string    `gorm:"not_null"`
	UserID uuid.UUID `gorm:"not_null;index"`
}
