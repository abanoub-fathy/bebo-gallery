package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service type contains all the types of the
// services used in our application
type Service struct {
	GalleryService
	UserService
}

// NewService is used to create service struct
// that will used to wrap all the services
// we used in our application inside it
func NewService(DB_URI string) (*Service, error) {
	// open db connection to be used in all services
	db, err := gorm.Open(postgres.Open(DB_URI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	service := &Service{
		GalleryService: NewGalleryGorm(db),
		UserService:    NewUserService(db),
	}

	return service, nil
}
