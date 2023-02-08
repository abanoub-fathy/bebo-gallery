package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Gallery is the container for images we will add
type Gallery struct {
	Base
	Title  string    `gorm:"not_null"`
	UserID uuid.UUID `gorm:"not_null;index"`
}

type GalleryService interface {
	GalleryDB
}

// GalleryDB has all methods needed to implemnt and
// use the Gallery database methods
type GalleryDB interface {
	// CreateGallery is used to create a new gallery into the DB
	CreateGallery(gallery *Gallery) error
}

// galleryGorm is the type that will implements the
// the GalleryDB for gorm
type galleryGorm struct {
	db *gorm.DB
}

// making sure that galleryGorm implemnts the GalleryDB
var _ GalleryDB = (*galleryGorm)(nil)

func (gg *galleryGorm) CreateGallery(gallery *Gallery) error {
	return gg.db.Create(&gallery).Error
}
